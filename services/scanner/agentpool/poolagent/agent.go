package poolagent

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-node/metrics"
	"github.com/forta-network/forta-node/nodeutils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/scanner"

	log "github.com/sirupsen/logrus"
)

// Constants
const (
	DefaultBufferSize             = 2000
	AgentTimeout                  = 30 * time.Second
	MaxFindings                   = 50
	DefaultAgentInitializeTimeout = 5 * time.Minute
)

// Agent receives blocks and transactions, and produces results.
type Agent struct {
	ctx    context.Context
	config config.AgentConfig

	txRequests          chan *TxRequest // never closed - deallocated when agent is discarded
	txResults           chan<- *scanner.TxResult
	blockRequests       chan *BlockRequest // never closed - deallocated when agent is discarded
	blockResults        chan<- *scanner.BlockResult
	combinationRequests chan *CombinationRequest // never closed - deallocated when agent is discarded
	combinationResults  chan<- *scanner.CombinationAlertResult

	errCounter *nodeutils.ErrorCounter
	msgClient  clients.MessageClient

	client         clients.AgentClient
	ready          chan struct{}
	readyOnce      sync.Once
	closed         chan struct{}
	closeOnce      sync.Once
	initialized    chan struct{}
	initializeOnce sync.Once

	mu sync.RWMutex
}

func (agent *Agent) AlertConfig() *protocol.AlertConfig {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	return agent.config.AlertConfig
}
func (agent *Agent) SetAlertConfig(cfg *protocol.AlertConfig) {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	agent.config.AlertConfig = cfg
}

func (agent *Agent) IsCombinerBot() bool {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	if agent.config.AlertConfig == nil {
		return false
	}

	return len(agent.config.AlertConfig.Subscriptions) > 0
}
func (agent *Agent) CombinerBotSubscriptions() []domain.CombinerBotSubscription {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	var subscriptions []domain.CombinerBotSubscription
	if !agent.IsCombinerBot() {
		return subscriptions
	}

	for _, subscription := range agent.AlertConfig().Subscriptions {
		subscriptions = append(
			subscriptions, domain.CombinerBotSubscription{
				Subscription: subscription,
				Subscriber: &domain.Subscriber{
					BotID:    agent.Config().ID,
					BotOwner: agent.Config().Owner,
					BotImage: agent.Config().Image,
				},
			},
		)
	}
	return subscriptions
}

// TxRequest contains the original request data and the encoded message.
type TxRequest struct {
	Original *protocol.EvaluateTxRequest
	Encoded  *grpc.PreparedMsg
}

// BlockRequest contains the original request data and the encoded message.
type BlockRequest struct {
	Original *protocol.EvaluateBlockRequest
	Encoded  *grpc.PreparedMsg
}

// CombinationRequest contains the original request data and the encoded message.
type CombinationRequest struct {
	Original *protocol.EvaluateAlertRequest
	Encoded  *grpc.PreparedMsg
}

// New creates a new agent.
func New(ctx context.Context, agentCfg config.AgentConfig, msgClient clients.MessageClient, txResults chan<- *scanner.TxResult, blockResults chan<- *scanner.BlockResult, alertResults chan<- *scanner.CombinationAlertResult) *Agent {
	return &Agent{
		ctx:                 ctx,
		config:              agentCfg,
		txRequests:          make(chan *TxRequest, DefaultBufferSize),
		txResults:           txResults,
		blockRequests:       make(chan *BlockRequest, DefaultBufferSize),
		blockResults:        blockResults,
		combinationRequests: make(chan *CombinationRequest, DefaultBufferSize),
		combinationResults:  alertResults,
		errCounter:          nodeutils.NewErrorCounter(3, isCriticalErr),
		msgClient:           msgClient,
		ready:               make(chan struct{}),
		closed:              make(chan struct{}),
		initialized:         make(chan struct{}),
	}
}

func isCriticalErr(err error) bool {
	return false
	// errStr := err.Error()
	// return strings.Contains(errStr, codes.DeadlineExceeded.String()) ||
	// 	strings.Contains(errStr, codes.Unavailable.String())
}

// LogStatus logs the status of the agent.
func (agent *Agent) LogStatus() {
	log.WithFields(log.Fields{
		"agent":       agent.config.ID,
		"blockBuffer": len(agent.blockRequests),
		"txBuffer":    len(agent.txRequests),
		"ready":       agent.IsReady(),
		"closed":      agent.IsClosed(),
	}).Debug("agent status")
}

// TxBufferIsFull tells if an agent input buffer is full.
func (agent *Agent) TxBufferIsFull() bool {
	return len(agent.txRequests) == DefaultBufferSize
}

// Config returns the agent config.
func (agent *Agent) Config() config.AgentConfig {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	return agent.config
}

// TxRequestCh returns the transaction request channel safely.
func (agent *Agent) TxRequestCh() chan<- *TxRequest {
	return agent.txRequests
}

// BlockRequestCh returns the block request channel safely.
func (agent *Agent) BlockRequestCh() chan<- *BlockRequest {
	return agent.blockRequests
}

// CombinationRequestCh returns the alert request channel safely.
func (agent *Agent) CombinationRequestCh() chan<- *CombinationRequest {
	return agent.combinationRequests
}

// Close implements io.Closer.
func (agent *Agent) Close() error {
	agent.closeOnce.Do(func() {
		close(agent.closed) // never close this anywhere else
		if agent.client != nil {
			agent.client.Close()
		}
	},
	)
	return nil
}

// SetReady sets the agent ready.
func (agent *Agent) SetReady() {
	agent.readyOnce.Do(
		func() {
			close(agent.ready) // never close this anywhere else
		},
	)
}

// SetInitialized sets the agent initialized.
func (agent *Agent) SetInitialized() {
	agent.initializeOnce.Do(
		func() {
			close(agent.initialized) // never close this anywhere else
		},
	)
}

// Ready returns the ready channel.
func (agent *Agent) Ready() <-chan struct{} {
	return agent.ready
}

// Initialized returns the initialized channel.
func (agent *Agent) Initialized() <-chan struct{} {
	return agent.initialized
}

// Closed returns the closed channel.
func (agent *Agent) Closed() <-chan struct{} {
	return agent.closed
}

// IsReady tells if the agent is ready.
func (agent *Agent) IsReady() bool {
	return isChanClosed(agent.ready)
}

// IsInitialized tells if the agent is initialized.
func (agent *Agent) IsInitialized() bool {
	return isChanClosed(agent.initialized)
}

// IsClosed tells if the agent is closed.
func (agent *Agent) IsClosed() bool {
	return isChanClosed(agent.closed)
}

func isChanClosed(ch chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

// SetClient sets the agent client for sending the requests.
func (agent *Agent) SetClient(agentClient clients.AgentClient) {
	agent.client = agentClient
}

// StartProcessing launches the goroutines to concurrently process incoming requests
// from request channels.
func (agent *Agent) StartProcessing() {
	go agent.processTransactions()
	go agent.processBlocks()
	go agent.processCombinationAlerts()
	return
}

func (agent *Agent) Initialize() {
	agentConfig := agent.Config()

	logger := log.WithFields(
		log.Fields{
			"agent": agentConfig.ID,
		},
	)

	// public bot.start metric to track bot starts/restarts.
	agent.msgClient.PublishProto(
		messaging.SubjectMetricAgent, &protocol.AgentMetricList{
			Metrics: []*protocol.AgentMetric{
				{
					AgentId:   agentConfig.ID,
					Timestamp: time.Now().Format(time.RFC3339),
					Name:      metrics.MetricStart,
					Value:     1,
				},
			},
		},
	)

	ctx, cancel := context.WithTimeout(agent.ctx, DefaultAgentInitializeTimeout)
	defer cancel()

	// invoke initialize method of the bot
	initializeResponse, err := agent.client.Initialize(ctx, &protocol.InitializeRequest{
		AgentId:   agentConfig.ID,
		ProxyHost: config.DockerJSONRPCProxyContainerName,
	})

	// it is not mandatory to implement a initialize method, safe to skip
	if status.Code(err) == codes.Unimplemented {
		logger.WithError(err).Info("Initialize() method not implemented in bot - safe to ignore")
		agent.SetInitialized()
		return
	}

	if err != nil {
		logger.WithError(err).Warn("bot initialization failed")
		_ = agent.Close()
		return
	}

	if err := validateInitializeResponse(initializeResponse); err != nil {
		logger.WithError(err).Warn("bot initialization validation failed")
		return
	}

	// pass new alert subscriptions to the agent pool
	if initializeResponse != nil && initializeResponse.AlertConfig != nil {
		agent.SetAlertConfig(initializeResponse.AlertConfig)
		agent.msgClient.Publish(messaging.SubjectAgentsAlertSubscribe, agent.CombinerBotSubscriptions())
	}

	logger.Info("bot initialization succeeded")
	agent.SetInitialized()
}

func validateInitializeResponse(response *protocol.InitializeResponse) error {
	if response == nil || response.AlertConfig == nil {
		return nil
	}

	for _, subscription := range response.AlertConfig.Subscriptions {
		if !utils.IsValidBotID(subscription.BotId) {
			return fmt.Errorf("invalid bot id :%s", subscription.BotId)
		}
	}

	return nil
}

func (agent *Agent) processTransactions() {
	lg := log.WithFields(
		log.Fields{
			"agent":     agent.Config().ID,
			"component": "agent",
			"evaluate":  "transaction",
		},
	)

	for request := range agent.txRequests {
		if exit := agent.processTransaction(lg, request); exit {
			return
		}
	}
}

func (agent *Agent) processTransaction(lg *log.Entry, request *TxRequest) (exit bool) {
	agentConfig := agent.Config()

	startTime := time.Now()
	if agent.IsClosed() {
		return true
	}

	ctx, cancel := context.WithTimeout(agent.ctx, AgentTimeout)
	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateTxResponse)

	requestTime := time.Now().UTC()
	err := agent.client.Invoke(ctx, agentgrpc.MethodEvaluateTx, request.Encoded, resp)
	responseTime := time.Now().UTC()
	cancel()
	if err == nil {
		// truncate findings
		if len(resp.Findings) > MaxFindings {
			dropped := len(resp.Findings) - MaxFindings
			droppedMetric := metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricFindingsDropped, float64(dropped))
			agent.msgClient.PublishProto(
				messaging.SubjectMetricAgent,
				&protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}},
			)
			resp.Findings = resp.Findings[:MaxFindings]
		}
		var duration time.Duration
		resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
		lg.WithField("duration", duration).Debugf("request successful")

		if resp.Metadata == nil {
			resp.Metadata = make(map[string]string)
		}
		resp.Metadata["imageHash"] = agentConfig.ImageHash()

		ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
		ts.BotRequest = requestTime
		ts.BotResponse = responseTime

		agent.txResults <- &scanner.TxResult{
			AgentConfig: agentConfig,
			Request:     request.Original,
			Response:    resp,
			Timestamps:  ts,
		}
		lg.WithField("duration", time.Since(startTime)).Debugf("sent results")

		return false
	}

	lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
	if agent.errCounter.TooManyErrs(err) {
		lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
		agent.Close()
		agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agentConfig})
		agent.msgClient.PublishProto(
			messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: []*protocol.AgentMetric{
					{
						AgentId:   agentConfig.ID,
						Timestamp: time.Now().Format(time.RFC3339),
						Name:      metrics.MetricStop,
						Value:     1,
					},
				},
			},
		)
		return true
	}

	return false
}

func (agent *Agent) processBlocks() {
	lg := log.WithFields(
		log.Fields{
			"agent":     agent.Config().ID,
			"component": "agent",
			"evaluate":  "block",
		},
	)

	for request := range agent.blockRequests {
		if exit := agent.processBlock(lg, request); exit {
			return
		}
	}
}

func (agent *Agent) processBlock(lg *log.Entry, request *BlockRequest) (exit bool) {
	agentConfig := agent.Config()

	startTime := time.Now()
	if agent.IsClosed() {
		return true
	}

	ctx, cancel := context.WithTimeout(agent.ctx, AgentTimeout)
	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateBlockResponse)
	requestTime := time.Now().UTC()
	err := agent.client.Invoke(ctx, agentgrpc.MethodEvaluateBlock, request.Encoded, resp)
	responseTime := time.Now().UTC()
	cancel()
	if err == nil {
		// truncate findings
		if len(resp.Findings) > MaxFindings {
			dropped := len(resp.Findings) - MaxFindings
			droppedMetric := metrics.CreateAgentMetric(
				agentConfig.ID, metrics.MetricFindingsDropped, float64(dropped),
			)
			agent.msgClient.PublishProto(
				messaging.SubjectMetricAgent,
				&protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}},
			)
			resp.Findings = resp.Findings[:MaxFindings]
		}
		var duration time.Duration
		resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
		lg.WithField("duration", duration).Debugf("request successful")

		if resp.Metadata == nil {
			resp.Metadata = make(map[string]string)
		}
		resp.Metadata["imageHash"] = agentConfig.ImageHash()

		ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
		ts.BotRequest = requestTime
		ts.BotResponse = responseTime

		agent.blockResults <- &scanner.BlockResult{
			AgentConfig: agentConfig,
			Request:     request.Original,
			Response:    resp,
			Timestamps:  ts,
		}
		lg.WithField("duration", time.Since(startTime)).Debugf("sent results")

		return false
	}

	lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
	if agent.errCounter.TooManyErrs(err) {
		lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
		agent.Close()
		agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agentConfig})

		return true
	}

	return false
}

func (agent *Agent) processCombinationAlerts() {
	lg := log.WithFields(
		log.Fields{
			"agent":     agent.Config().ID,
			"component": "agent",
			"evaluate":  "combination",
		},
	)

	for request := range agent.combinationRequests {
		if exit := agent.processCombinationAlert(lg, request); exit {
			return
		}
	}
}

func (agent *Agent) processCombinationAlert(lg *log.Entry, request *CombinationRequest) bool {
	agentConfig := agent.Config()

	startTime := time.Now()
	if agent.IsClosed() {
		return true
	}

	ctx, cancel := context.WithTimeout(agent.ctx, AgentTimeout)
	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateAlertResponse)
	requestTime := time.Now().UTC()
	err := agent.client.Invoke(ctx, agentgrpc.MethodEvaluateAlert, request.Encoded, resp)
	responseTime := time.Now().UTC()
	cancel()

	if err != nil {
		lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
		if agent.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
			agent.Close()
			agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agentConfig})

			return true
		}
	}

	// validate response
	if vErr := validateEvaluateAlertResponse(resp); vErr != nil {
		lg.WithField(
			"request", request.Original.RequestId,
		).WithError(vErr).Error("evaluate combination response validation failed")

		return false
	}

	// truncate findings
	if len(resp.Findings) > MaxFindings {
		dropped := len(resp.Findings) - MaxFindings
		droppedMetric := metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricFindingsDropped, float64(dropped))
		agent.msgClient.PublishProto(
			messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}},
		)
		resp.Findings = resp.Findings[:MaxFindings]
	}

	var duration time.Duration
	resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
	lg.WithField("duration", duration).Debugf("request successful")

	if resp.Metadata == nil {
		resp.Metadata = make(map[string]string)
	}

	resp.Metadata["imageHash"] = agentConfig.ImageHash()

	ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
	ts.BotRequest = requestTime
	ts.BotResponse = responseTime

	agent.combinationResults <- &scanner.CombinationAlertResult{
		AgentConfig: agentConfig,
		Request:     request.Original,
		Response:    resp,
		Timestamps:  ts,
	}

	lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
	return false
}

func validateEvaluateAlertResponse(resp *protocol.EvaluateAlertResponse) (err error) {
	if resp == nil {
		return fmt.Errorf("nil response")
	}

	for _, finding := range resp.Findings {
		if err = validateFinding(finding); err != nil {
			return err
		}
	}

	return nil
}

func validateFinding(finding *protocol.Finding) error {
	if finding == nil {
		return fmt.Errorf("nil finding")
	}
	for _, alert := range finding.RelatedAlerts {
		if !checkValidKeccak256(alert) {
			return fmt.Errorf("bad related alert string: %s", alert)
		}
	}
	for _, address := range finding.Addresses {
		if !common.IsHexAddress(address) {
			return fmt.Errorf("bad address string: %s", address)
		}
	}

	return nil
}

var _regexKeccak256 = regexp.MustCompile("^0x[a-f0-9]{64}$")

func checkValidKeccak256(hash string) bool {
	return _regexKeccak256.Match([]byte(hash))
}

func calculateResponseTime(startTime *time.Time) (timestamp string, latencyMs uint32, duration time.Duration) {
	now := time.Now().UTC()
	duration = now.Sub(*startTime)
	return now.Format(time.RFC3339), uint32(duration.Milliseconds()), duration
}

// ShouldProcessBlock tells if the agent should process block.
func (agent *Agent) ShouldProcessBlock(blockNumberHex string) bool {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	blockNumber, _ := hexutil.DecodeUint64(blockNumberHex)
	var isAtLeastStartBlock bool
	if agent.config.StartBlock != nil {
		isAtLeastStartBlock = blockNumber >= *agent.config.StartBlock
	} else {
		isAtLeastStartBlock = true
	}

	var isAtMostStopBlock bool
	if agent.config.StopBlock != nil {
		isAtMostStopBlock = blockNumber <= *agent.config.StopBlock
	} else {
		isAtMostStopBlock = true
	}

	var isOnThisShard bool
	// if sharded, block % shards must be equal to shard id
	if agent.IsSharded() {
		isOnThisShard = uint(blockNumber)%agent.config.ShardConfig.Shards == agent.config.ShardConfig.ShardID
	} else {
		isOnThisShard = true
	}

	return isAtLeastStartBlock && isAtMostStopBlock && isOnThisShard
}

func (agent *Agent) ShouldProcessAlert(event *protocol.AlertEvent) bool {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	if agent.config.AlertConfig == nil {
		return false
	}

	// handle sharding
	alertCreatedAt, err := time.Parse(time.RFC3339Nano, event.Alert.CreatedAt)
	if err != nil {
		log.WithFields(
			log.Fields{
				"alertHash": event.Alert.Hash,
				"createdAt": event.Alert.CreatedAt,
				"botId":     agent.config.ID,
			},
		).Warn("failed to parse created at for sharding calculation")

		return false
	}

	var isOnThisShard bool
	if agent.IsSharded() {
		isOnThisShard = uint(alertCreatedAt.Unix())%agent.config.ShardConfig.Shards == agent.config.ShardConfig.ShardID
	} else {
		isOnThisShard = true
	}

	return isOnThisShard
}

func (agent *Agent) UpdateConfig(cfg config.AgentConfig) {
	agent.mu.Lock()
	defer agent.mu.Unlock()

	agent.config.ShardConfig = cfg.ShardConfig
	agent.config.Manifest = cfg.Manifest
}

func (agent *Agent) IsSharded() bool {
	agent.mu.RLock()
	defer agent.mu.RUnlock()

	return agent.config.ShardConfig != nil && agent.config.ShardConfig.Shards > 1
}
