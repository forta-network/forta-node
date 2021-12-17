package publisher

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/goccy/go-json"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/contracts"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services/publisher/testalerts"
	"github.com/forta-protocol/forta-node/utils"
)

const (
	defaultInterval        = time.Second * 15
	defaultBatchLimit      = 500
	defaultBatchBufferSize = 100
)

// Publisher receives, collects and publishes alerts.
type Publisher struct {
	protocol.UnimplementedPublisherNodeServer
	ctx               context.Context
	cfg               PublisherConfig
	contract          AlertsContract
	ipfs              IPFS
	testAlertLogger   TestAlertLogger
	metricsAggregator *AgentMetricsAggregator
	messageClient     *messaging.Client

	initialize    sync.Once
	port          int
	skipEmpty     bool
	skipPublish   bool
	batchInterval time.Duration
	batchLimit    int
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

type PublisherConfig struct {
	Port            int
	ChainID         int
	Key             *keystore.Key
	PublisherConfig config.PublisherConfig
}

func (pub *Publisher) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	pub.notifCh <- req
	return &protocol.NotifyResponse{}, nil
}

func (pub *Publisher) publishNextBatch(batch *protocol.SignedAlertBatch) error {
	// flush only if we are publishing so we can make the best use of aggregated metrics
	if _, skip := pub.shouldSkipPublishing(batch); !skip {
		batch.Data.Metrics = pub.metricsAggregator.TryFlush()
	}

	signature, err := security.SignProtoMessage(pub.cfg.Key, batch)
	if err != nil {
		return fmt.Errorf("failed to sign alert batch: %v", err)
	}
	batch.Signature = signature

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(batch); err != nil {
		return fmt.Errorf("failed to encode the signed alert: %v", err)
	}
	log.Tracef("alert payload: %s", string(buf.Bytes()))

	// save with blank tx hash for now
	if pub.skipPublish {
		log.Infof("alert batch: blockStart=%d, blockEnd=%d, alertCount=%d, maxSeverity=%s", batch.Data.BlockStart, batch.Data.BlockEnd, batch.Data.AlertCount, batch.Data.MaxSeverity.String())
		log.Info("skipping batch, because skipPublish is enabled")
		return nil
	}

	// if we should really skip this batch due to other reasons, then we just leave it in the db with blank tx hash
	if reason, skip := pub.shouldSkipPublishing(batch); skip {
		log.WithField("reason", reason).Info("skipping batch")
		return nil
	}

	cid, err := pub.ipfs.Add(&buf, ipfsapi.Pin(true))
	if err != nil {
		return fmt.Errorf("failed to store alert data to ipfs: %v", err)
	}

	logger := log.WithFields(
		log.Fields{
			"blockStart":  batch.Data.BlockStart,
			"blockEnd":    batch.Data.BlockEnd,
			"alertCount":  batch.Data.AlertCount,
			"maxSeverity": batch.Data.MaxSeverity.String(),
			"ref":         cid,
			"metrics":     len(batch.Data.Metrics),
		},
	)

	tx, err := pub.contract.AddAlertBatch(
		big.NewInt(0).SetUint64(batch.Data.ChainId),
		big.NewInt(0).SetUint64(batch.Data.BlockStart),
		big.NewInt(0).SetUint64(batch.Data.BlockEnd),
		big.NewInt(int64(batch.Data.AlertCount)),
		big.NewInt(0).SetUint64(uint64(batch.Data.MaxSeverity)),
		cid,
	)
	if err != nil {
		logger.WithError(err).Error("alert while sending batch")
		return fmt.Errorf("failed to send the alert tx: %v", err)
	}

	logger.WithFields(
		log.Fields{
			"tx": tx.Hash().Hex(),
		},
	).Info("alert batch")

	return nil
}

func (pub *Publisher) shouldSkipPublishing(batch *protocol.SignedAlertBatch) (string, bool) {
	return "because there are no alerts and skipEmpty is enabled",
		pub.skipEmpty && batch.Data.AlertCount == uint32(0)
}

func (pub *Publisher) listenForMetrics() {
	pub.messageClient.Subscribe(messaging.SubjectMetricAgent, messaging.AgentMetricHandler(pub.metricsAggregator.AddAgentMetrics))
}

func (pub *Publisher) publishBatches() {
	for batch := range pub.batchCh {
		if err := pub.publishNextBatch(batch); err != nil {
			log.Errorf("failed to publish alert batch: %v", err)
			time.Sleep(time.Second * 3)
		}
		time.Sleep(time.Millisecond * 200)
	}
}

func (pub *Publisher) prepareBatches() {
	for {
		pub.prepareLatestBatch()
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

func (pub *Publisher) prepareLatestBatch() {
	batch := &BatchData{Data: &protocol.AlertBatch{ChainId: uint64(pub.cfg.ChainID)}}

	timeoutCh := time.After(pub.batchInterval)

	var done bool
	var i int
	for i < pub.batchLimit {
		select {
		case notif := <-pub.notifCh:
			alert := notif.SignedAlert
			hasAlert := alert != nil
			if hasAlert {
				log.Debugf("alert: %s", alert.Alert.Id)
			}

			if hasAlert && notif.SignedAlert.Alert.Agent.IsTest {
				if pub.cfg.PublisherConfig.TestAlerts.Disable {
					continue
				}
				if err := pub.testAlertLogger.LogTestAlert(pub.ctx, notif.SignedAlert); err != nil {
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

	pub.batchCh <- (*protocol.SignedAlertBatch)(batch)
}

// on connection of first agent, start publishing batches (no agents = no batches)
func (pub *Publisher) handleReady(cfgs messaging.AgentPayload) error {
	pub.initialize.Do(func() {
		go pub.prepareBatches()
		go pub.publishBatches()
		go pub.listenForMetrics()
	})
	return nil
}

func (pub *Publisher) Start() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", pub.port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	protocol.RegisterPublisherNodeServer(grpcServer, pub)

	pub.messageClient.Subscribe(messaging.SubjectAgentsStatusAttached, messaging.AgentsHandler(pub.handleReady))

	return grpcServer.Serve(lis)
}

func (pub *Publisher) Stop() error {
	log.Infof("Stopping %s", pub.Name())
	return nil
}

func (pub *Publisher) Name() string {
	return "Publisher"
}

func NewPublisher(ctx context.Context, mc *messaging.Client, cfg PublisherConfig) (*Publisher, error) {
	rpcClient, err := rpc.Dial(cfg.PublisherConfig.JsonRpc.Url)

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

	txOpts.GasPrice = big.NewInt(cfg.PublisherConfig.GasPriceGwei * params.GWei)
	txOpts.GasLimit = cfg.PublisherConfig.GasLimit

	contract, err := contracts.NewAlertsTransactor(common.HexToAddress(cfg.PublisherConfig.ContractAddress), ethereum.NewContractBackend(rpcClient))
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

	return &Publisher{
		ctx:               ctx,
		cfg:               cfg,
		contract:          ats,
		ipfs:              ipfsClient,
		testAlertLogger:   testAlertLogger,
		metricsAggregator: NewMetricsAggregator(),
		messageClient:     mc,

		port:          cfg.Port,
		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		skipPublish:   cfg.PublisherConfig.SkipPublish,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.SignedAlertBatch, defaultBatchBufferSize),
	}, nil
}
