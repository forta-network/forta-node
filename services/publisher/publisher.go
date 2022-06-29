package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/clients/webhook"
	"github.com/forta-network/forta-core-go/clients/webhook/client/operations"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/ipfs"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/protocol/transform"
	"github.com/forta-network/forta-core-go/release"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/alertapi"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/store"
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
	ipfs              ipfs.Client
	metricsAggregator *AgentMetricsAggregator
	messageClient     *messaging.Client
	alertClient       clients.AlertAPIClient
	webhookClient     webhook.AlertWebhookClient

	batchRefStore    store.StringStore
	lastReceiptStore store.StringStore

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
		var flushed bool
		batch.Metrics, flushed = pub.metricsAggregator.TryFlush()
		if flushed {
			log.Debug("flushed metrics")
			pub.lastMetricsFlush.Set()
		} else {
			log.Debug("not flushing metrics yet")
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
	lastBatchRef, err := pub.batchRefStore.Get()
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
		log.WithFields(log.Fields{
			"blockStart":  batch.BlockStart,
			"blockEnd":    batch.BlockEnd,
			"alertCount":  batch.AlertCount,
			"maxSeverity": batch.MaxSeverity.String(),
		}).Info(reason)
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

	if pub.cfg.Config.LocalModeConfig.Enable {
		scannerJwt, err := security.CreateScannerJWT(pub.cfg.Key, map[string]interface{}{
			"localMode": "true",
		})
		alertBatch := transform.ToWebhookAlertBatch(batch)
		if !pub.cfg.Config.LocalModeConfig.IncludeMetrics {
			log.Debug("excluding metrics due to local mode config")
			alertBatch.Metrics = nil
		}
		_, err = pub.webhookClient.SendAlerts(&operations.SendAlertsParams{
			Context:       context.Background(),
			Payload:       alertBatch,
			Authorization: utils.StringPtr(fmt.Sprintf("Bearer %s", scannerJwt)),
		})
		if err != nil {
			log.WithError(err).Error("failed to send private alerts")
			return err
		}
		if alertBatch != nil {
			log.WithFields(log.Fields{
				"alertCount":   len(alertBatch.Alerts),
				"metricsCount": len(alertBatch.Metrics),
			}).Info("successfully sent private alerts")
		}
		return nil
	}

	cid, err := pub.ipfs.CalculateFileHash(buf.Bytes())
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

	var lastReceipt string
	lr, err := pub.lastReceiptStore.Get()
	if err == nil {
		lastReceipt = lr
	}

	signedBatchSummary, err := security.SignBatchSummary(pub.cfg.Key, &protocol.BatchSummary{
		Batch:            cid,
		ChainId:          batch.ChainId,
		BlockStart:       batch.BlockStart,
		BlockEnd:         batch.BlockEnd,
		AlertCount:       batch.AlertCount,
		ScannerVersion:   batch.ScannerVersion,
		PreviousReceipt:  lastReceipt,
		LatestBlockInput: batch.LatestBlockInput,
		Timestamp:        time.Now().UTC().Format(time.RFC3339),
	})
	if err != nil {
		logger.WithError(err).Error("failed to sign batch summary")
		return err
	}

	scannerJwt, err := security.CreateScannerJWT(pub.cfg.Key, map[string]interface{}{
		"batch": cid,
	})

	if err != nil {
		logger.WithError(err).Error("failed to sign cid")
		return err
	}
	resp, err := pub.alertClient.PostBatch(&domain.AlertBatchRequest{
		Scanner:            pub.cfg.Key.Address.Hex(),
		ChainID:            int64(batch.ChainId),
		BlockStart:         int64(batch.BlockStart),
		BlockEnd:           int64(batch.BlockEnd),
		AlertCount:         int64(batch.AlertCount),
		MaxSeverity:        int64(batch.MaxSeverity),
		Ref:                cid,
		SignedBatch:        signedBatch,
		SignedBatchSummary: signedBatchSummary,
	}, scannerJwt)

	if err != nil {
		logger.WithError(err).Error("alert while sending batch")
		return fmt.Errorf("failed to send the alert tx: %v", err)
	}

	//TODO: after receipts are returned, make it non-optional
	if resp.SignedReceipt != nil {
		// store off receipt id
		if err := pub.lastReceiptStore.Put(resp.ReceiptID); err != nil {
			logger.WithError(err).Error("failed to marshal receipt")
			return err
		}
		logger = logger.WithFields(log.Fields{
			"receiptId": resp.ReceiptID,
		})

		// if for some reason receipt can't marshal, log and move on
		b, err := json.Marshal(resp.SignedReceipt)
		if err != nil {
			logger.WithError(err).Error("failed to marshal receipt (not saving receipt)")
			return nil
		}
		logger = logger.WithFields(log.Fields{
			"receipt": string(b),
		})
	}

	logger.Info("alert batch")

	return nil
}

func (pub *Publisher) shouldSkipPublishing(batch *protocol.AlertBatch) (string, bool) {
	if batch.AlertCount > 0 {
		return "", false
	}
	// after this line, alert count is considered as zero but let's check metrics
	localModeMetricsOnly := pub.cfg.Config.LocalModeConfig.Enable && len(batch.Metrics) > 0
	if localModeMetricsOnly {
		return "", false
	}
	const defaultReason = "because there are no alerts"
	if pub.cfg.Config.LocalModeConfig.Enable {
		return defaultReason + " or metrics and local mode is enabled", true
	}
	if pub.skipEmpty {
		return defaultReason + " and skipEmpty is enabled", true
	}
	// should not skip: there must be a combining reason along with the zero count
	return "", false
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
		}
		time.Sleep(time.Millisecond * 20)
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
				log.WithField("alertId", alert.Alert.Id).Debug("publisher received alert")
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
	go pub.prepareBatches()
	go pub.publishBatches()
	pub.registerMessageHandlers()
	return nil
}

func (pub *Publisher) Stop() error {
	cfg := pub.cfg.Config
	if cfg.LocalModeConfig.Enable {
		timeoutSeconds := cfg.LocalModeConfig.RuntimeLimits.StopTimeoutSeconds
		log.WithField("timeout", fmt.Sprintf("%ds", timeoutSeconds)).Info("waiting for scanning to finish")
		time.Sleep(time.Duration(timeoutSeconds) * time.Second)
		log.WithField("timeout", fmt.Sprintf("%ds", timeoutSeconds)).Info("done waiting scanning to finish")
	}
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

func NewPublisher(ctx context.Context, cfg config.Config) (*Publisher, error) {
	mc := messaging.NewClient("metrics", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	releaseInfoStr := os.Getenv(config.EnvReleaseInfo)
	var releaseSummary *release.ReleaseSummary
	if len(releaseInfoStr) > 0 {
		releaseInfo := release.ReleaseInfoFromString(releaseInfoStr)
		releaseSummary = release.MakeSummaryFromReleaseInfo(releaseInfo)
	}

	apiClient := alertapi.NewClient(cfg.Publish.APIURL)

	return initPublisher(ctx, mc, apiClient, PublisherConfig{
		ChainID:         cfg.ChainID,
		Key:             key,
		PublisherConfig: cfg.Publish,
		ReleaseSummary:  releaseSummary,
		Config:          cfg,
	})
}

func initPublisher(ctx context.Context, mc *messaging.Client, alertClient clients.AlertAPIClient, cfg PublisherConfig) (*Publisher, error) {
	ipfsClient, err := ipfs.NewClient(fmt.Sprintf("http://%s:5001", config.DockerIpfsContainerName))
	if err != nil {
		return nil, err
	}

	batchInterval := defaultInterval
	if cfg.PublisherConfig.Batch.IntervalSeconds != nil {
		batchInterval = (time.Duration)(*cfg.PublisherConfig.Batch.IntervalSeconds) * time.Second
	}

	batchLimit := defaultBatchLimit
	if cfg.PublisherConfig.Batch.MaxAlerts != nil {
		batchLimit = *cfg.PublisherConfig.Batch.MaxAlerts
	}

	var webhookClient webhook.AlertWebhookClient
	if cfg.Config.LocalModeConfig.Enable {
		dest := cfg.Config.LocalModeConfig.WebhookURL
		webhookClient, err = webhook.NewAlertWebhookClient(dest)
		if err != nil {
			return nil, fmt.Errorf("invalid private alert webhook url: %s", dest)
		}
	}

	return &Publisher{
		ctx:               ctx,
		cfg:               cfg,
		ipfs:              ipfsClient,
		metricsAggregator: NewMetricsAggregator(time.Duration(*cfg.PublisherConfig.Batch.MetricsBucketIntervalSeconds) * time.Second),
		messageClient:     mc,
		alertClient:       alertClient,
		webhookClient:     webhookClient,
		batchRefStore:     store.NewFileStringStore(path.Join(cfg.Config.FortaDir, ".last-batch")),
		lastReceiptStore:  store.NewFileStringStore(path.Join(cfg.Config.FortaDir, ".last-receipt")),

		skipEmpty:     cfg.PublisherConfig.Batch.SkipEmpty,
		skipPublish:   cfg.PublisherConfig.SkipPublish,
		batchInterval: batchInterval,
		batchLimit:    batchLimit,
		notifCh:       make(chan *protocol.NotifyRequest, defaultBatchLimit),
		batchCh:       make(chan *protocol.AlertBatch, defaultBatchBufferSize),
	}, nil
}
