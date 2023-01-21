package poolagent

import (
	"context"
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/ethereum/go-ethereum/common"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/utils"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-node/metrics"
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
	MaxFindings                   = 10
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

	errCounter *errorCounter
	msgClient  clients.MessageClient

	client    clients.AgentClient
	ready     chan struct{}
	readyOnce sync.Once
	closed    chan struct{}
	closeOnce sync.Once

	// initialization fields
	initWait         sync.WaitGroup

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
func New(
	ctx context.Context, agentCfg config.AgentConfig, msgClient clients.MessageClient,
	txResults chan<- *scanner.TxResult, blockResults chan<- *scanner.BlockResult,
	alertResults chan<- *scanner.CombinationAlertResult,
) *Agent {
	return &Agent{
		ctx:                 ctx,
		config:              agentCfg,
		txRequests:          make(chan *TxRequest, DefaultBufferSize),
		txResults:           txResults,
		blockRequests:       make(chan *BlockRequest, DefaultBufferSize),
		blockResults:        blockResults,
		combinationRequests: make(chan *CombinationRequest, DefaultBufferSize),
		combinationResults:  alertResults,
		errCounter:          NewErrorCounter(3, isCriticalErr),
		msgClient:           msgClient,
		ready:               make(chan struct{}),
		closed:              make(chan struct{}),
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
	})
	return nil
}

// SetReady sets the agent ready.
func (agent *Agent) SetReady() {
	agent.readyOnce.Do(func() {
		close(agent.ready) // never close this anywhere else
	})
}

// Ready returns the ready channel.
func (agent *Agent) Ready() <-chan struct{} {
	return agent.ready
}

// Closed returns the closed channel.
func (agent *Agent) Closed() <-chan struct{} {
	return agent.closed
}

// IsReady tells if the agent is ready.
func (agent *Agent) IsReady() bool {
	return isChanClosed(agent.ready)
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
	agent.initWait.Add(1)

	go func() {
		err := agent.startInitWorker(agent.ctx)
		if err != nil {
			log.WithError(err).WithField("botId", agent.config.ID).Warn("failed to initialize bot")
		}
	}()
	go agent.processTransactions()
	go agent.processBlocks()
	go agent.processCombinationAlerts()
}

func (agent *Agent) startInitWorker(ctx context.Context) error {
	bExp := &backoff.ExponentialBackOff{
		InitialInterval:     time.Second * 30,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         time.Minute * 10,
		MaxElapsedTime:      time.Hour,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}

	bo := backoff.WithContext(bExp, ctx)
	err := backoff.Retry(
		func() error {
			initCtx, cancel := context.WithTimeout(agent.ctx, DefaultAgentInitializeTimeout)
			defer cancel()
			return agent.initialize(initCtx)
		}, bo,
	)

	return fmt.Errorf("agent.initialize() backoff failed: %v", err)
}

func (agent *Agent) initialize(ctx context.Context) error {
	logger := log.WithFields(
		log.Fields{
			"botId": agent.config.ID,
		},
	)

	initializeResponse, err := agent.client.Initialize(
		ctx, &protocol.InitializeRequest{
			AgentId:   agent.config.ID,
			ProxyHost: config.DockerJSONRPCProxyContainerName,
		},
	)

	if status.Code(err) == codes.Unimplemented {
		logger.WithError(err).Info("initialize() method not implemented in bot - safe to ignore")
		return nil
	}

	if err != nil {
		logger.WithError(err).Warn("bot initialization failed")
		return err
	}

	if err := validateInitializeResponse(initializeResponse); err != nil {
		logger.WithError(err).Warn("bot initialization validation failed")
		return err
	}

	// pass new alert subscriptions to pool
	if initializeResponse != nil {
		agent.SetAlertConfig(initializeResponse.AlertConfig)
		for _, subscription := range initializeResponse.AlertConfig.Subscriptions {
			agent.msgClient.Publish(
				messaging.SubjectAgentsAlertSubscribe, messaging.SubscriptionPayload{
					messaging.
					CombinerBotSubscription{Subscription: subscription},
				},
			)
		}
	}

	logger.Info("bot initialization succeeded")

	agent.initWait.Done()

	return nil
}

func validateInitializeResponse(response *protocol.InitializeResponse) error {
	if response == nil || response.AlertConfig == nil{
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
	lg := log.WithFields(log.Fields{
		"agent":     agent.config.ID,
		"component": "agent",
		"evaluate":  "transaction",
	})

	agent.initWait.Wait()

	for request := range agent.txRequests {
		startTime := time.Now()
		if agent.IsClosed() {
			return
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
				droppedMetric := metrics.CreateAgentMetric(agent.config.ID, metrics.MetricFindingsDropped, float64(dropped))
				agent.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}})
				resp.Findings = resp.Findings[:MaxFindings]
			}
			var duration time.Duration
			resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
			lg.WithField("duration", duration).Debugf("request successful")

			if resp.Metadata == nil {
				resp.Metadata = make(map[string]string)
			}
			resp.Metadata["imageHash"] = agent.config.ImageHash()

			ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
			ts.BotRequest = requestTime
			ts.BotResponse = responseTime

			agent.txResults <- &scanner.TxResult{
				AgentConfig: agent.config,
				Request:     request.Original,
				Response:    resp,
				Timestamps:  ts,
			}
			lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
			continue
		}
		lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
		if agent.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
			agent.Close()
			agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent.config})
			agent.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{
				Metrics: []*protocol.AgentMetric{
					{
						AgentId:   agent.config.ID,
						Timestamp: time.Now().Format(time.RFC3339),
						Name:      metrics.MetricStop,
						Value:     1,
					},
				},
			})
			return
		}
	}
}

func (agent *Agent) processBlocks() {
	lg := log.WithFields(log.Fields{
		"agent":     agent.config.ID,
		"component": "agent",
		"evaluate":  "block",
	})

	agent.initWait.Wait()

	for request := range agent.blockRequests {
		startTime := time.Now()
		if agent.IsClosed() {
			return
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
				droppedMetric := metrics.CreateAgentMetric(agent.config.ID, metrics.MetricFindingsDropped, float64(dropped))
				agent.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}})
				resp.Findings = resp.Findings[:MaxFindings]
			}
			var duration time.Duration
			resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
			lg.WithField("duration", duration).Debugf("request successful")

			if resp.Metadata == nil {
				resp.Metadata = make(map[string]string)
			}
			resp.Metadata["imageHash"] = agent.config.ImageHash()

			ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
			ts.BotRequest = requestTime
			ts.BotResponse = responseTime

			agent.blockResults <- &scanner.BlockResult{
				AgentConfig: agent.config,
				Request:     request.Original,
				Response:    resp,
				Timestamps:  ts,
			}
			lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
			continue
		}
		lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking agent")
		if agent.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down agent")
			agent.Close()
			agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent.config})
			return
		}
	}
}

func (agent *Agent) processCombinationAlerts() {
	lg := log.WithFields(
		log.Fields{
			"agent":     agent.config.ID,
			"component": "agent",
			"evaluate":  "combination",
		},
	)

	agent.initWait.Wait()

	for request := range agent.combinationRequests {
		startTime := time.Now()
		if agent.IsClosed() {
			return
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
				agent.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent.config})
				return
			}
		}

		// validate response
		if vErr := validateEvaluateAlertResponse(resp); vErr != nil {
			lg.WithField("request", request.Original.RequestId).WithError(vErr).Error("evaluate combination response validation failed")
			continue
		}

		// truncate findings
		if len(resp.Findings) > MaxFindings {
			dropped := len(resp.Findings) - MaxFindings
			droppedMetric := metrics.CreateAgentMetric(agent.config.ID, metrics.MetricFindingsDropped, float64(dropped))
			agent.msgClient.PublishProto(messaging.SubjectMetricAgent, &protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{droppedMetric}})
			resp.Findings = resp.Findings[:MaxFindings]
		}

		var duration time.Duration
		resp.Timestamp, resp.LatencyMs, duration = calculateResponseTime(&startTime)
		lg.WithField("duration", duration).Debugf("request successful")

		if resp.Metadata == nil {
			resp.Metadata = make(map[string]string)
		}

		resp.Metadata["imageHash"] = agent.config.ImageHash()

		ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
		ts.BotRequest = requestTime
		ts.BotResponse = responseTime

		agent.combinationResults <- &scanner.CombinationAlertResult{
			AgentConfig: agent.config,
			Request:     request.Original,
			Response:    resp,
			Timestamps:  ts,
		}

		lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
		continue
	}
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

	return isAtLeastStartBlock && isAtMostStopBlock
}

func (agent *Agent) ShouldProcessAlert(event *protocol.AlertEvent) bool {
	if agent.config.AlertConfig == nil {
		return false
	}
	for _, subscription := range agent.config.AlertConfig.Subscriptions {
		// bot is subscribed to the bot id
		subscribedToBot := subscription.BotId == "" || subscription.BotId == event.Alert.Source.Bot.Id
		// bot is subscribed to the alert id
		subscribedToAlert := subscription.AlertId == "" || subscription.AlertId == event.Alert.AlertId

		if subscribedToBot && subscribedToAlert {
			return true
		}
	}

	return false
}
