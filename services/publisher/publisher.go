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

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/domain"
	"github.com/forta-protocol/forta-core-go/protocol"
	"github.com/forta-protocol/forta-core-go/release"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/services/publisher/testalerts"
	"github.com/forta-protocol/forta-node/store"
	"github.com/goccy/go-json"
	ipfsapi "github.com/ipfs/go-ipfs-api"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
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
	alertClient       clients.AlertAPIClient
	batchRefStore     store.BatchRefStore

	server *grpc.Server

	initialize    sync.Once
	skipEmpty     bool
	skipPublish   bool
	batchInterval time.Duration
	batchLimit    int
	latestChainID uint64
	notifCh       chan *protocol.NotifyRequest
	batchCh       chan *protocol.AlertBatch

	lastBatchPublish    health.TimeTracker
	lastBatchSkip       health.TimeTracker
	lastBatchSkipReason health.MessageTracker
	lastBatchPublishErr health.ErrorTracker
	lastMetricsFlush    health.TimeTracker

	latestBlockInput   uint64
	latestBlockInputMu sync.RWMutex
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
	ChainID         int
	Key             *keystore.Key
	PublisherConfig config.PublisherConfig
	ReleaseSummary  *release.ReleaseSummary
	Config          config.Config
}

func (pub *Publisher) Notify(ctx context.Context, req *protocol.NotifyRequest) (*protocol.NotifyResponse, error) {
	pub.notifCh <- req
	return &protocol.NotifyResponse{}, nil
}

func (pub *Publisher) publishNextBatch(batch *protocol.AlertBatch) error {
	// flush only if we are publishing so we can make the best use of aggregated metrics
	if _, skip := pub.shouldSkipPublishing(batch); !skip {
		batch.Metrics = pub.metricsAggregator.TryFlush()
		if len(batch.Metrics) > 0 {
			pub.lastMetricsFlush.Set()
		}
	}

	// add release info if it's available
	if pub.cfg.ReleaseSummary != nil {
		batch.ScannerVersion = &protocol.ScannerVersion{
			Commit:  pub.cfg.ReleaseSummary.Commit,
			Ipfs:    pub.cfg.ReleaseSummary.IPFS,
			Version: pub.cfg.ReleaseSummary.Version,
		}
	}
	lastBatchRef, err := pub.batchRefStore.GetLast()
	if err == nil {
		batch.Parent = lastBatchRef
	}

	// use the latest block input from scanner, fall back to latest block number from the batch
	pub.latestBlockInputMu.RLock()
	batch.LatestBlockInput = pub.latestBlockInput
	pub.latestBlockInputMu.RUnlock()
	if batch.LatestBlockInput == 0 {
		batch.LatestBlockInput = batch.BlockEnd
	}

	signedBatch, err := security.SignBatch(pub.cfg.Key, batch)
	if err != nil {
		return fmt.Errorf("failed to build envelope: %v", err)
	}

	var buf bytes.Buffer
	if err = json.NewEncoder(&buf).Encode(signedBatch); err != nil {
		return fmt.Errorf("failed to encode the signed alert: %v", err)
	}
	log.Tracef("alert payload: %s", string(buf.Bytes()))

	if pub.skipPublish {
		const reason = "skipping batch, because skipPublish is enabled"
		log.Infof("alert batch: blockStart=%d, blockEnd=%d, alertCount=%d, maxSeverity=%s", batch.BlockStart, batch.BlockEnd, batch.AlertCount, batch.MaxSeverity.String())
		log.Info(reason)
		pub.lastBatchSkip.Set()
		pub.lastBatchSkipReason.Set(reason)
		return nil
	}

	if reason, skip := pub.shouldSkipPublishing(batch); skip {
		log.WithField("reason", reason).Info("skipping batch")
		pub.lastBatchSkip.Set()
		pub.lastBatchSkipReason.Set(reason)
		return nil
	}

	cid, err := pub.ipfs.Add(&buf, ipfsapi.OnlyHash(true))
	if err != nil {
		return fmt.Errorf("failed to store alert data to ipfs: %v", err)
	}
	if err := pub.batchRefStore.Put(cid); err != nil {
		return fmt.Errorf("failed to write last batch ref: %v", err)
	}

	logger := log.WithFields(
		log.Fields{
			"blockStart":  batch.BlockStart,
			"blockEnd":    batch.BlockEnd,
			"alertCount":  batch.AlertCount,
			"maxSeverity": batch.MaxSeverity.String(),
			"ref":         cid,
			"metrics":     len(batch.Metrics),
		},
	)

	scannerJwt, err := security.CreateScannerJWT(pub.cfg.Key, map[string]interface{}{
		"batch": cid,
	})

	if err != nil {
		logger.WithError(err).Error("failed to sign cid")
		return err
	}
	err = pub.alertClient.PostBatch(&domain.AlertBatchRequest{
		Scanner:     pub.cfg.Key.Address.Hex(),
		ChainID:     int64(batch.ChainId),
		BlockStart:  int64(batch.BlockStart),
		BlockEnd:    int64(batch.BlockEnd),
		AlertCount:  int64(batch.AlertCount),
		MaxSeverity: int64(batch.MaxSeverity),
		Ref:         cid,
		SignedBatch: signedBatch,
	}, scannerJwt)

	if err != nil {
		logger.WithError(err).Error("alert while sending batch")
		return fmt.Errorf("failed to send the alert tx: %v", err)
	}

	logger.Info("alert batch")

	return nil
}

func (pub *Publisher) shouldSkipPublishing(batch *protocol.AlertBatch) (string, bool) {
	return "because there are no alerts and skipEmpty is enabled",
		pub.skipEmpty && batch.AlertCount == uint32(0)
}

func (pub *Publisher) registerMessageHandlers() {
	pub.messageClient.Subscribe(messaging.SubjectMetricAgent, messaging.AgentMetricHandler(pub.metricsAggregator.AddAgentMetrics))
	pub.messageClient.Subscribe(messaging.SubjectScannerBlock, messaging.ScannerHandler(pub.handleScannerBlock))
}

func (pub *Publisher) handleScannerBlock(payload messaging.ScannerPayload) error {
	pub.latestBlockInputMu.Lock()
	defer pub.latestBlockInputMu.Unlock()

	logger := log.WithFields(log.Fields{
		"newLatestBlockInput":  payload.LatestBlockInput,
		"prevLatestBlockInput": pub.latestBlockInput,
	})
	if payload.LatestBlockInput < pub.latestBlockInput {
		logger.Warn("skipping scanner update (lower than previous)")
		return nil
	}
	logger.Info("received scanner update")
	pub.latestBlockInput = payload.LatestBlockInput
	return nil
}

func (pub *Publisher) publishBatches() {
	for batch := range pub.batchCh {
		err := pub.publishNextBatch(batch)
		pub.lastBatchPublish.Set()
		pub.lastBatchPublishErr.Set(err)
		if err != nil {
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

// BatchData is a parent wrapper that contains all batch info.
type BatchData protocol.AlertBatch

func (bd *BatchData) GetPrivateAlerts(notif *protocol.NotifyRequest) *protocol.AgentAlerts {
	for _, a := range bd.PrivateAlerts {
		if a.AgentManifest == notif.AgentInfo.Manifest {
			return a
		}
	}
	res := &protocol.AgentAlerts{
		AgentManifest: notif.AgentInfo.Manifest,
	}

	bd.PrivateAlerts = append(bd.PrivateAlerts, res)
	return res
}

// AppendAlert adds the alert to the relevant list.
func (bd *BatchData) AppendAlert(notif *protocol.NotifyRequest) {
	isBlockAlert := notif.EvalBlockRequest != nil

	var isPrivate bool

	if notif.SignedAlert != nil && notif.SignedAlert.Alert != nil && notif.SignedAlert.Alert.Finding != nil {
		// default at per-finding level
		isPrivate = notif.SignedAlert.Alert.Finding.Private

		// if public, let a private override at response-level win
		if !isPrivate {
			if notif.EvalBlockResponse != nil {
				isPrivate = notif.EvalBlockResponse.Private
			} else if notif.EvalTxResponse != nil {
				isPrivate = notif.EvalTxResponse.Private
			}
		}
	}

	hasAlert := notif.SignedAlert != nil

	var agentAlerts *protocol.AgentAlerts
	if isPrivate {
		if hasAlert {
			agentAlerts = bd.GetPrivateAlerts(notif)
		}
	} else if isBlockAlert {
		blockNum := hexutil.MustDecodeUint64(notif.EvalBlockRequest.Event.BlockNumber)
		bd.AddBatchAgent(notif.AgentInfo, blockNum, "")
		blockRes := bd.GetBlockResults(notif.EvalBlockRequest.Event.BlockHash, blockNum, notif.EvalBlockRequest.Event.Block.Timestamp)
		if hasAlert {
			agentAlerts = (*BlockResults)(blockRes).GetAgentAlerts(notif.AgentInfo)
		}
	} else {
		blockNum := hexutil.MustDecodeUint64(notif.EvalTxRequest.Event.Block.BlockNumber)
		bd.AddBatchAgent(notif.AgentInfo, blockNum, notif.EvalTxRequest.Event.Receipt.TransactionHash)
		blockRes := bd.GetBlockResults(notif.EvalTxRequest.Event.Block.BlockHash, blockNum, notif.EvalTxRequest.Event.Block.BlockTimestamp)
		if hasAlert {
			txRes := (*BlockResults)(blockRes).GetTransactionResults(notif.EvalTxRequest.Event)
			agentAlerts = (*TransactionResults)(txRes).GetAgentAlerts(notif.AgentInfo)
		}
	}

	if agentAlerts == nil {
		return
	}

	agentAlerts.Alerts = append(agentAlerts.Alerts, notif.SignedAlert)
	bd.AlertCount++
}

// AddBatchAgent includes the agent info in the batch so we know that this agent really
// processed a specific block or a tx hash.
func (bd *BatchData) AddBatchAgent(agent *protocol.AgentInfo, blockNumber uint64, txHash string) {
	var batchAgent *protocol.BatchAgent
	for _, ba := range bd.Agents {
		if ba.Info.Manifest == agent.Manifest {
			batchAgent = ba
			break
		}
	}
	if batchAgent == nil {
		batchAgent = &protocol.BatchAgent{
			Info: agent,
		}
		bd.Agents = append(bd.Agents, batchAgent)
	}
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
func (bd *BatchData) GetBlockResults(blockHash string, blockNumber uint64, blockTimestamp string) *protocol.BlockResults {
	for _, blockRes := range bd.Results {
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
	bd.Results = append(bd.Results, br)
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
	batch := (*BatchData)(&protocol.AlertBatch{ChainId: uint64(pub.cfg.ChainID)})

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
			if batch.BlockStart == 0 || (batch.BlockStart > 0 && notifBlockNum < batch.BlockStart) {
				batch.BlockStart = notifBlockNum
			}
			if batch.BlockEnd == 0 || (batch.BlockEnd > 0 && notifBlockNum > batch.BlockEnd) {
				batch.BlockEnd = notifBlockNum
			}

			if hasAlert && alert.Alert.Finding.Severity > batch.MaxSeverity {
				batch.MaxSeverity = alert.Alert.Finding.Severity
			}

			batch.AppendAlert(notif)

		case <-timeoutCh:
			done = true
		}

		if done {
			break
		}
	}

	pub.batchCh <- (*protocol.AlertBatch)(batch)
}

func (pub *Publisher) Start() error {
	lis, err := net.Listen("tcp", "0.0.0.0:8770")
	if err != nil {
		return err
	}
	pub.server = grpc.NewServer()
	protocol.RegisterPublisherNodeServer(pub.server, pub)

	go pub.prepareBatches()
	go pub.publishBatches()
	pub.registerMessageHandlers()

	utils.GoGrpcServe(pub.server, lis)
	return nil
}

func (pub *Publisher) Stop() error {
	log.Infof("Stopping %s", pub.Name())
	if pub.server != nil {
		pub.server.Stop()
	}
	return nil
}

func (pub *Publisher) Name() string {
	return "publisher"
}

// Health implements the health.Reporter interface.
func (pub *Publisher) Health() health.Reports {
	return health.Reports{
		pub.lastBatchPublish.GetReport("event.batch-publish.time"),
		pub.lastBatchPublishErr.GetReport("event.batch-publish.error"),
		&health.Report{
			Name:    "event.batch-skip.time",
			Status:  health.StatusInfo,
			Details: pub.lastBatchSkip.String(),
		},
		pub.lastBatchSkipReason.GetReport("event.batch-skip.reason"),
		pub.lastMetricsFlush.GetReport("event.metrics-flush.time"),
	}
}

func NewPublisher(ctx context.Context, mc *messaging.Client, alertClient clients.AlertAPIClient, cfg PublisherConfig) (*Publisher, error) {

	var ipfsClient IPFS
	switch {
	case cfg.PublisherConfig.SkipPublish:
		// use nil IPFS client

	case len(cfg.PublisherConfig.IPFS.Username) > 0 && len(cfg.PublisherConfig.IPFS.Password) > 0:
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.IPFS.APIURL, &http.Client{
			Transport: utils.NewBasicAuthTransport(cfg.PublisherConfig.IPFS.Username, cfg.PublisherConfig.IPFS.Password),
		})

	default:
		ipfsClient = ipfsapi.NewShellWithClient(cfg.PublisherConfig.IPFS.APIURL, http.DefaultClient)
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
		ipfs:              ipfsClient,
		testAlertLogger:   testAlertLogger,
		metricsAggregator: NewMetricsAggregator(),
		messageClient:     mc,
		alertClient:       alertClient,
		batchRefStore:     store.NewBatchRefStore(cfg.Config.FortaDir),

		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		skipPublish:   cfg.PublisherConfig.SkipPublish,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.AlertBatch, defaultBatchBufferSize),
	}, nil
}
