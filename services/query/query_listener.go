package query

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"time"

	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/contracts"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/security"
	"OpenZeppelin/fortify-node/store"
	"OpenZeppelin/fortify-node/utils"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultInterval   = time.Second * 15
	defaultBatchLimit = 500
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx       context.Context
	store     store.AlertStore
	cfg       AlertListenerConfig
	contract  AlertsContract
	ipfs      IPFS
	ethClient EthClient

	skipEmpty     bool
	batchInterval time.Duration
	batchLimit    int
	latestBlock   uint64
	latestChainID uint64
	notifCh       chan *protocol.NotifyRequest
}

// EthClient interacts with the Ethereum API.
type EthClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

// AlertsContract stores alerts.
type AlertsContract interface {
	AddAlertBatch(_chainId *big.Int, _blockStart *big.Int, _blockEnd *big.Int, _alertCount *big.Int, _maxSeverity *big.Int, _ref string) (*types.Transaction, error)
}

// IPFS interacts with an IPFS node/gateway.
type IPFS interface {
	Add(r io.Reader, options ...ipfsapi.AddOpts) (string, error)
}

type AlertListenerConfig struct {
	Port            int
	ChainID         int
	Key             *keystore.Key
	PublisherConfig config.PublisherConfig
}

func (al *AlertListener) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	al.notifCh <- req
	return &protocol.NotifyResponse{}, nil
}

func (al *AlertListener) publishAlerts() {
	ticker := time.NewTicker(al.batchInterval)
	var batchErr error
	for {
		if batchErr != nil {
			log.Errorf("failed to publish alert batch: %v", batchErr)
			// Sleep
			ticker.Reset(al.batchInterval)
			<-ticker.C
		}

		batch := al.getLatestBatch()

		signature, err := security.SignProtoMessage(al.cfg.Key, (*protocol.SignedAlertBatch)(batch))
		if err != nil {
			batchErr = fmt.Errorf("failed to sign alert batch: %v", err)
			continue
		}
		batch.Signature = signature

		var buf bytes.Buffer
		if err = json.NewEncoder(&buf).Encode(batch); err != nil {
			batchErr = fmt.Errorf("failed to encode the signed alert: %v", err)
			continue
		}
		log.Infof("new alert batch: %s", string(buf.Bytes()))
		cid, err := al.ipfs.Add(&buf, ipfsapi.Pin(true))
		if err != nil {
			batchErr = fmt.Errorf("failed to store alert data to ipfs: %v", err)
			continue
		}
		log.Infof("stored the batch to ipfs: %s", cid)

		tx, err := al.contract.AddAlertBatch(
			big.NewInt(0).SetUint64(batch.Data.ChainId),
			big.NewInt(0).SetUint64(batch.Data.BlockStart),
			big.NewInt(0).SetUint64(batch.Data.BlockEnd),
			big.NewInt(int64(batch.Data.AlertCount)),
			big.NewInt(0).SetUint64(uint64(batch.Data.MaxSeverity)),
			cid,
		)
		if err != nil {
			batchErr = fmt.Errorf("failed to send the alert tx: %v", err)
			continue
		}

		// Store all block and transaction alerts.
		al.storeBatchWithTxHash(batch.Data, tx.Hash().Hex())

		batchErr = nil
		<-ticker.C
	}
}

func (al *AlertListener) storeBatchWithTxHash(batch *protocol.AlertBatch, txHash string) {
	for _, result := range batch.Results {
		for _, blockRes := range result.Results {
			for _, alert := range blockRes.Alerts {
				al.storeAlertWithTxHash(alert, txHash)
			}
		}
		for _, txs := range result.Transactions {
			for _, txRes := range txs.Results {
				for _, alert := range txRes.Alerts {
					al.storeAlertWithTxHash(alert, txHash)
				}
			}
		}
	}
}

func (al *AlertListener) storeAlertWithTxHash(alert *protocol.SignedAlert, txHash string) {
	alert.PublishedWithTx = txHash
	if err := al.store.AddAlert(alert); err != nil {
		log.Errorf("failed to store the alert: %v", err)
	}
}

// TransactionResults contains the results for a transaction.
type TransactionResults protocol.TransactionResults

// BlockResults contains the results for a block.
type BlockResults protocol.BlockResults

// AlertBatch contains the actual batch data.
type AlertBatch protocol.AlertBatch

// BatchData is a parent wrapper that contains all batch info.
type BatchData protocol.SignedAlertBatch

// AppendAlert adds the alert to the relevant list.
func (bd *BatchData) AppendAlert(notif *protocol.NotifyRequest) {
	alertBatch := (*AlertBatch)(bd.Data)
	isBlockAlert := notif.EvalBlockRequest != nil
	hasAlert := notif.SignedAlert != nil

	var agentAlerts *protocol.AgentAlerts
	if isBlockAlert {
		blockNum := hexutil.MustDecodeUint64(notif.EvalBlockRequest.Event.BlockNumber)
		blockRes := alertBatch.GetBlockResults(notif.EvalBlockRequest.Event.BlockHash, blockNum)
		if hasAlert {
			agentAlerts = (*BlockResults)(blockRes).GetAgentAlerts(notif.SignedAlert.Alert.Agent)
		}
	} else {
		blockNum := hexutil.MustDecodeUint64(notif.EvalTxRequest.Event.Block.BlockNumber)
		blockRes := alertBatch.GetBlockResults(notif.EvalTxRequest.Event.Block.BlockHash, blockNum)
		txRes := (*BlockResults)(blockRes).GetTransactionResults(notif.EvalTxRequest.Event)
		if hasAlert {
			agentAlerts = (*TransactionResults)(txRes).GetAgentAlerts(notif.SignedAlert.Alert.Agent)
		}
	}

	if !hasAlert {
		return
	}

	agentAlerts.Alerts = append(agentAlerts.Alerts, notif.SignedAlert)
	bd.Data.AlertCount++
}

// GetBlockResults returns an existing or a new aggregation object for the block.
func (ab *AlertBatch) GetBlockResults(blockHash string, blockNumber uint64) *protocol.BlockResults {
	for _, blockRes := range ab.Results {
		if blockRes.Block.BlockNumber == blockNumber {
			return blockRes
		}
	}
	br := &protocol.BlockResults{
		Block: &protocol.Block{
			BlockHash:   blockHash,
			BlockNumber: blockNumber,
		},
	}
	ab.Results = append(ab.Results, br)
	return br
}

// GetTransactionResults returns an existing or a new aggregation object for the transaction.
func (br *BlockResults) GetTransactionResults(tx *protocol.TransactionEvent) *protocol.TransactionResults {
	for _, txRes := range br.Transactions {
		if txRes.Transaction.Transaction.Hash == tx.Transaction.Hash {
			return txRes
		}
	}
	tr := &protocol.TransactionResults{
		Transaction: tx,
	}
	br.Transactions = append(br.Transactions, tr)
	return tr
}

// GetAgentAlerts returns an existing or a new aggregation object for the agent alerts.
func (br *BlockResults) GetAgentAlerts(agent *protocol.AgentInfo) *protocol.AgentAlerts {
	for _, agentAlerts := range br.Results {
		if agentAlerts.Agent.Name == agent.Name {
			return agentAlerts
		}
	}
	aa := &protocol.AgentAlerts{
		Agent: agent,
	}
	br.Results = append(br.Results, aa)
	return aa
}

// GetAgentAlerts returns an existing or a new aggregation object for the agent alerts.
func (tr *TransactionResults) GetAgentAlerts(agent *protocol.AgentInfo) *protocol.AgentAlerts {
	for _, agentAlerts := range tr.Results {
		if agentAlerts.Agent.Name == agent.Name {
			return agentAlerts
		}
	}
	aa := &protocol.AgentAlerts{
		Agent: agent,
	}
	tr.Results = append(tr.Results, aa)
	return aa
}

func (al *AlertListener) getLatestBatch() (batch *BatchData) {
	batch = &BatchData{Data: &protocol.AlertBatch{ChainId: uint64(al.cfg.ChainID)}}

	var done bool
	for i := 0; i < defaultBatchLimit; i++ {
		select {
		case notif := <-al.notifCh:
			alert := notif.SignedAlert
			hasAlert := alert != nil
			if hasAlert {
				log.Infof("alert: %s", alert.Alert.Id)
			}

			var blockNum string
			if notif.EvalBlockRequest != nil {
				blockNum = notif.EvalBlockRequest.Event.BlockNumber
			} else {
				blockNum = notif.EvalTxRequest.Event.Block.BlockNumber
			}

			notifBlockNum, err := hexutil.DecodeUint64(blockNum)
			if err != nil {
				log.Errorf("failed to parse alert notif block number: %v", err)
				continue
			}
			if batch.Data.BlockStart == 0 || (batch.Data.BlockStart > 0 && notifBlockNum < batch.Data.BlockStart) {
				batch.Data.BlockStart = notifBlockNum
			}
			if batch.Data.BlockEnd == 0 || (batch.Data.BlockEnd > 0 && notifBlockNum > batch.Data.BlockEnd) {
				batch.Data.BlockEnd = notifBlockNum
			}

			if hasAlert && alert.Alert.Finding.Severity > batch.Data.MaxSeverity {
				batch.Data.MaxSeverity = alert.Alert.Finding.Severity
			}

			batch.AppendAlert(notif)

		default:
			done = true // If we don't receive anymore notifs
		}
		if done {
			break
		}
	}

	if batch.Data.BlockStart == 0 {
		latestBlock, err := al.ethClient.BlockNumber(al.ctx)
		if err != nil {
			log.Errorf("failed to get the latest block for batch: %v", err)
			return
		}
		batch.Data.BlockStart = al.latestBlock
		batch.Data.BlockEnd = latestBlock
		al.latestBlock = latestBlock
	}

	return
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)

	go al.publishAlerts()

	return grpcServer.Serve(lis)
}

func (al *AlertListener) Stop() error {
	log.Infof("Stopping %s", al.Name())
	return nil
}

func (al *AlertListener) Name() string {
	return "AlertListener"
}

func NewAlertListener(ctx context.Context, store store.AlertStore, cfg AlertListenerConfig) (*AlertListener, error) {
	rpcClient, err := rpc.Dial(cfg.PublisherConfig.Ethereum.JsonRpcUrl)
	if err != nil {
		return nil, err
	}
	ethClient := ethclient.NewClient(rpcClient)
	txOpts := bind.NewKeyedTransactor(cfg.Key.PrivateKey)
	contract, err := contracts.NewAlertsTransactor(common.HexToAddress(cfg.PublisherConfig.ContractAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the alerts contract: %v", err)
	}
	ats := &contracts.AlertsTransactorSession{
		Contract:     contract,
		TransactOpts: *txOpts,
	}

	var ipfsClient IPFS
	if len(cfg.PublisherConfig.IPFS.Username) > 0 && len(cfg.PublisherConfig.IPFS.Password) > 0 {
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.IPFS.GatewayURL, &http.Client{
			Transport: utils.NewBasicAuthTransport(cfg.PublisherConfig.IPFS.Username, cfg.PublisherConfig.IPFS.Password),
		})
	} else {
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.IPFS.GatewayURL, http.DefaultClient)
	}

	latestBlock, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get the latest block: %v", err)
	}

	batchInterval := defaultInterval
	if cfg.PublisherConfig.Batch.IntervalSeconds != nil {
		batchInterval = (time.Duration)(*cfg.PublisherConfig.Batch.IntervalSeconds) * time.Second
	}

	batchLimit := defaultBatchLimit
	if cfg.PublisherConfig.Batch.MaxAlerts != nil {
		batchLimit = *cfg.PublisherConfig.Batch.MaxAlerts
	}

	return &AlertListener{
		ctx:       ctx,
		store:     store,
		cfg:       cfg,
		contract:  ats,
		ipfs:      ipfsClient,
		ethClient: ethClient,

		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		latestBlock:   latestBlock,
		notifCh:       make(chan *protocol.NotifyRequest, batchLimit),
	}, nil
}
