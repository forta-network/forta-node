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

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/contracts"
	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/security"
	"github.com/forta-network/forta-node/services/query/testalerts"
	"github.com/forta-network/forta-node/store"
	"github.com/forta-network/forta-node/utils"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

const (
	defaultInterval        = time.Second * 15
	defaultBatchLimit      = 500
	defaultBatchBufferSize = 100
)

// AlertListener allows retrieval of alerts from the database
type AlertListener struct {
	protocol.UnimplementedQueryNodeServer
	ctx             context.Context
	store           store.AlertStore
	cfg             AlertListenerConfig
	contract        AlertsContract
	ipfs            IPFS
	ethClient       EthClient
	testAlertLogger TestAlertLogger

	port          int
	skipEmpty     bool
	skipPublish   bool
	batchInterval time.Duration
	batchLimit    int
	latestBlock   uint64
	latestChainID uint64
	notifCh       chan *protocol.NotifyRequest
	batchCh       chan *protocol.SignedAlertBatch
}

// TestAlertLogger logs the test alerts.
type TestAlertLogger interface {
	LogTestAlert(context.Context, *protocol.SignedAlert) error
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

func (al *AlertListener) publishNextBatch(batch *protocol.SignedAlertBatch) error {
	signature, err := security.SignProtoMessage(al.cfg.Key, batch)
	if err != nil {
		return fmt.Errorf("failed to sign alert batch: %v", err)
	}
	batch.Signature = signature

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(batch); err != nil {
		return fmt.Errorf("failed to encode the signed alert: %v", err)
	}
	log.Debugf("alert payload: %s", string(buf.Bytes()))

	al.storeBatchWithTxHash(batch.Data, "")
	if al.skipPublish {
		log.Infof("alert batch: blockStart=%d, blockEnd=%d, alertCount=%d, maxSeverity=%s", batch.Data.BlockStart, batch.Data.BlockEnd, batch.Data.AlertCount, batch.Data.MaxSeverity.String())
		log.Info("skipping batch, because skipPublish is enabled")
		return nil
	}

	// if no alerts, and skipEmpty is true, then save with blank txHash
	if al.skipEmpty && batch.Data.AlertCount == uint32(0) {
		log.Info("skipping batch, because there are no alerts and skipEmpty is enabled")
		return nil
	}

	cid, err := al.ipfs.Add(&buf, ipfsapi.Pin(true))
	if err != nil {
		return fmt.Errorf("failed to store alert data to ipfs: %v", err)
	}
	log.Infof("alert batch: blockStart=%d, blockEnd=%d, alertCount=%d, maxSeverity=%s, ref=%s", batch.Data.BlockStart, batch.Data.BlockEnd, batch.Data.AlertCount, batch.Data.MaxSeverity.String(), cid)

	tx, err := al.contract.AddAlertBatch(
		big.NewInt(0).SetUint64(batch.Data.ChainId),
		big.NewInt(0).SetUint64(batch.Data.BlockStart),
		big.NewInt(0).SetUint64(batch.Data.BlockEnd),
		big.NewInt(int64(batch.Data.AlertCount)),
		big.NewInt(0).SetUint64(uint64(batch.Data.MaxSeverity)),
		cid,
	)
	if err != nil {
		return fmt.Errorf("failed to send the alert tx: %v", err)
	}

	// Store all block and transaction alerts.
	al.storeBatchWithTxHash(batch.Data, tx.Hash().Hex())
	return nil
}

func (al *AlertListener) publishBatches() {
	for batch := range al.batchCh {
		if err := al.publishNextBatch(batch); err != nil {
			log.Errorf("failed to publish alert batch: %v", err)
			time.Sleep(time.Second * 3)
		}
		time.Sleep(time.Second * 200)
	}
}

func (al *AlertListener) prepareBatches() {
	for {
		al.prepareLatestBatch()
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
		alertBatch.AddBatchAgent(notif.AgentInfo, blockNum, "")
		if hasAlert {
			blockRes := alertBatch.GetBlockResults(notif.EvalBlockRequest.Event.BlockHash, blockNum, notif.EvalBlockRequest.Event.Block.Timestamp)
			agentAlerts = (*BlockResults)(blockRes).GetAgentAlerts(notif.AgentInfo)
		}
	} else {
		blockNum := hexutil.MustDecodeUint64(notif.EvalTxRequest.Event.Block.BlockNumber)
		alertBatch.AddBatchAgent(notif.AgentInfo, blockNum, notif.EvalTxRequest.Event.Receipt.TransactionHash)
		if hasAlert {
			blockRes := alertBatch.GetBlockResults(notif.EvalTxRequest.Event.Block.BlockHash, blockNum, notif.EvalTxRequest.Event.Block.BlockTimestamp)
			txRes := (*BlockResults)(blockRes).GetTransactionResults(notif.EvalTxRequest.Event)
			agentAlerts = (*TransactionResults)(txRes).GetAgentAlerts(notif.AgentInfo)
		}
	}

	if agentAlerts == nil {
		return
	}

	agentAlerts.Alerts = append(agentAlerts.Alerts, notif.SignedAlert)
	bd.Data.AlertCount++
}

// AddBatchAgent includes the agent info in the batch so we know that this agent really
// processed a specific block or a tx hash.
func (ab *AlertBatch) AddBatchAgent(agent *protocol.AgentInfo, blockNumber uint64, txHash string) {
	var batchAgent *protocol.BatchAgent
	for _, ba := range ab.Agents {
		if ba.Info.Manifest == agent.Manifest {
			batchAgent = ba
			break
		}
	}
	if batchAgent == nil {
		batchAgent = &protocol.BatchAgent{
			Info: agent,
		}
		ab.Agents = append(ab.Agents, batchAgent)
	}
	// There should always be a block number.
	if blockNumber == 0 {
		log.Error("zero block number while adding batch agent")
		return
	}
	var alreadyAddedBlockNum bool
	for _, addedBlockNum := range batchAgent.Blocks {
		if addedBlockNum == blockNumber {
			alreadyAddedBlockNum = true
			break
		}
	}
	if !alreadyAddedBlockNum {
		batchAgent.Blocks = append(batchAgent.Blocks, blockNumber)
	}
	if len(txHash) > 0 {
		batchAgent.Transactions = append(batchAgent.Transactions, txHash)
	}
}

// GetBlockResults returns an existing or a new aggregation object for the block.
func (ab *AlertBatch) GetBlockResults(blockHash string, blockNumber uint64, blockTimestamp string) *protocol.BlockResults {
	for _, blockRes := range ab.Results {
		if blockRes.Block.BlockNumber == blockNumber {
			return blockRes
		}
	}
	br := &protocol.BlockResults{
		Block: &protocol.Block{
			BlockHash:      blockHash,
			BlockNumber:    blockNumber,
			BlockTimestamp: blockTimestamp,
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
		if agentAlerts.AgentManifest == agent.Manifest {
			return agentAlerts
		}
	}
	aa := &protocol.AgentAlerts{
		AgentManifest: agent.Manifest,
	}
	br.Results = append(br.Results, aa)
	return aa
}

// GetAgentAlerts returns an existing or a new aggregation object for the agent alerts.
func (tr *TransactionResults) GetAgentAlerts(agent *protocol.AgentInfo) *protocol.AgentAlerts {
	for _, agentAlerts := range tr.Results {
		if agentAlerts.AgentManifest == agent.Manifest {
			return agentAlerts
		}
	}
	aa := &protocol.AgentAlerts{
		AgentManifest: agent.Manifest,
	}
	tr.Results = append(tr.Results, aa)
	return aa
}

func (al *AlertListener) prepareLatestBatch() {
	batch := &BatchData{Data: &protocol.AlertBatch{ChainId: uint64(al.cfg.ChainID)}}

	timeoutCh := time.After(defaultInterval)

	var done bool
	var i int
	for i < al.batchLimit {
		select {
		case notif := <-al.notifCh:
			alert := notif.SignedAlert
			hasAlert := alert != nil
			if hasAlert {
				log.Debugf("alert: %s", alert.Alert.Id)
			}

			if hasAlert && notif.SignedAlert.Alert.Agent.IsTest {
				if al.cfg.PublisherConfig.TestAlerts.Disable {
					continue
				}
				if err := al.testAlertLogger.LogTestAlert(al.ctx, notif.SignedAlert); err != nil {
					log.Warnf("failed to log test alert: %v", err)
				}
				continue
			}

			// Notifications with empty alerts shouldn't be taken into account while limiting the batch.
			// Otherwise, we create too many batches very quickly.
			if hasAlert {
				i++
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

		case <-timeoutCh:
			done = true
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

	al.batchCh <- (*protocol.SignedAlertBatch)(batch)
}

func (al *AlertListener) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", al.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterQueryNodeServer(grpcServer, al)

	go al.prepareBatches()
	go al.publishBatches()

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
	chainID, err := ethClient.ChainID(ctx)
	if err != nil {
		log.Errorf("could not determine scanner registry chain ID: %s", err.Error())
		return nil, err
	}

	txOpts, err := bind.NewKeyedTransactorWithChainID(cfg.Key.PrivateKey, chainID)
	if err != nil {
		log.Errorf("error while creating keyed transactor for listener: %s", err.Error())
		return nil, err
	}

	contract, err := contracts.NewAlertsTransactor(common.HexToAddress(cfg.PublisherConfig.ContractAddress), ethClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize the alerts contract: %v", err)
	}

	ats := &contracts.AlertsTransactorSession{
		Contract:     contract,
		TransactOpts: *txOpts,
	}

	var ipfsClient IPFS
	switch {
	case cfg.PublisherConfig.SkipPublish:
		// use nil IPFS client

	case len(cfg.PublisherConfig.IPFS.Username) > 0 && len(cfg.PublisherConfig.IPFS.Password) > 0:
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.IPFS.GatewayURL, &http.Client{
			Transport: utils.NewBasicAuthTransport(cfg.PublisherConfig.IPFS.Username, cfg.PublisherConfig.IPFS.Password),
		})

	default:
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

	var testAlertLogger TestAlertLogger
	if !cfg.PublisherConfig.TestAlerts.Disable {
		testAlertLogger = testalerts.NewLogger(cfg.PublisherConfig.TestAlerts.WebhookURL)
	}

	return &AlertListener{
		ctx:             ctx,
		store:           store,
		cfg:             cfg,
		contract:        ats,
		ipfs:            ipfsClient,
		ethClient:       ethClient,
		testAlertLogger: testAlertLogger,

		port:          cfg.Port,
		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		skipPublish:   cfg.PublisherConfig.SkipPublish,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		latestBlock:   latestBlock,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.SignedAlertBatch, defaultBatchBufferSize),
	}, nil
}
