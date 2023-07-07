package botio

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/nodeutils"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	"github.com/forta-network/forta-node/services/components/metrics"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// BotClient represents a detection bot that is being communicated to and managed.
type BotClient interface {
	Config() config.AgentConfig
	SetConfig(config.AgentConfig)

	Initialized() <-chan struct{}
	IsInitialized() bool
	Closed() <-chan struct{}
	IsClosed() bool

	TxBufferIsFull() bool

	Initialize()
	StartProcessing()

	ShouldProcessBlock(blockNumberHex string) bool
	ShouldProcessAlert(event *protocol.AlertEvent) bool

	TxRequestCh() chan<- *botreq.TxRequest
	BlockRequestCh() chan<- *botreq.BlockRequest
	CombinationRequestCh() chan<- *botreq.CombinationRequest

	LogStatus()

	CombinerBotSubscriptions() []domain.CombinerBotSubscription

	io.Closer
}

// Constants
const (
	DefaultBufferSize          = 2000
	RequestTimeout             = 30 * time.Second
	MaxFindings                = 50
	DefaultInitializeTimeout   = 5 * time.Minute
	DefaultHealthCheckInterval = time.Second * 30
)

// botClient receives blocks and transactions, and produces results.
type botClient struct {
	ctx               context.Context
	ctxCancel         func()
	configUnsafe      config.AgentConfig
	alertConfigUnsafe protocol.AlertConfig

	txRequests          chan *botreq.TxRequest          // never closed - deallocated when bot is discarded
	blockRequests       chan *botreq.BlockRequest       // never closed - deallocated when bot is discarded
	combinationRequests chan *botreq.CombinationRequest // never closed - deallocated when bot is discarded

	resultChannels botreq.SendOnlyChannels

	errCounter       *nodeutils.ErrorCounter
	msgClient        clients.MessageClient
	lifecycleMetrics metrics.Lifecycle

	dialer       agentgrpc.BotDialer
	clientUnsafe agentgrpc.Client

	initialized     chan struct{}
	initializedOnce sync.Once

	closeOnce sync.Once

	mu sync.RWMutex
}

var _ BotClient = &botClient{}

func (bot *botClient) isCombinerBot() bool {
	return len(bot.AlertConfig().Subscriptions) > 0
}

func (bot *botClient) CombinerBotSubscriptions() []domain.CombinerBotSubscription {
	return MakeCombinerBotSubscriptions(bot.AlertConfig().Subscriptions, bot.Config())
}

// MakeCombinerBotSubscriptions makes combiner bot subscriptions from given alert config subscriptions.
func MakeCombinerBotSubscriptions(
	alertSubs []*protocol.CombinerBotSubscription,
	botConfig config.AgentConfig,
) (subscriptions []domain.CombinerBotSubscription) {
	for _, subscription := range alertSubs {
		subscriptions = append(
			subscriptions, domain.CombinerBotSubscription{
				Subscription: subscription,
				Subscriber: &domain.Subscriber{
					BotID:    botConfig.ID,
					BotOwner: botConfig.Owner,
					BotImage: botConfig.Image,
				},
			},
		)
	}
	return
}

// NewBotClient creates a new bot client.
func NewBotClient(
	ctx context.Context, botCfg config.AgentConfig,
	msgClient clients.MessageClient, lifecycleMetrics metrics.Lifecycle, botDialer agentgrpc.BotDialer,
	resultChannels botreq.SendOnlyChannels,
) *botClient {
	botCtx, botCtxCancel := context.WithCancel(ctx)
	return &botClient{
		ctx:                 botCtx,
		ctxCancel:           botCtxCancel,
		configUnsafe:        botCfg,
		txRequests:          make(chan *botreq.TxRequest, DefaultBufferSize),
		blockRequests:       make(chan *botreq.BlockRequest, DefaultBufferSize),
		combinationRequests: make(chan *botreq.CombinationRequest, DefaultBufferSize),
		resultChannels:      resultChannels,
		errCounter:          nodeutils.NewErrorCounter(3, isCriticalErr),
		msgClient:           msgClient,
		lifecycleMetrics:    lifecycleMetrics,
		dialer:              botDialer,
		initialized:         make(chan struct{}),
	}
}

func isCriticalErr(err error) bool {
	return false
	// errStr := err.Error()
	// return strings.Contains(errStr, codes.DeadlineExceeded.String()) ||
	// 	strings.Contains(errStr, codes.Unavailable.String())
}

// LogStatus logs the status of the bot.
func (bot *botClient) LogStatus() {
	log.WithFields(log.Fields{
		"bot":         bot.Config().ID,
		"blockBuffer": len(bot.blockRequests),
		"txBuffer":    len(bot.txRequests),
		"initialized": bot.IsInitialized(),
		"closed":      bot.IsClosed(),
	}).Debug("bot status")
}

// TxBufferIsFull tells if an bot input buffer is full.
func (bot *botClient) TxBufferIsFull() bool {
	return len(bot.txRequests) == DefaultBufferSize
}

// SetConfig sets the bot config.
func (bot *botClient) SetConfig(botConfig config.AgentConfig) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	bot.configUnsafe = botConfig
}

// Config returns the bot config.
func (bot *botClient) Config() config.AgentConfig {
	bot.mu.RLock()
	defer bot.mu.RUnlock()

	return bot.configUnsafe
}

// AlertConfig returns the alert config.
func (bot *botClient) AlertConfig() protocol.AlertConfig {
	bot.mu.RLock()
	defer bot.mu.RUnlock()

	return bot.alertConfigUnsafe
}

// SetAlertConfig sets the bot config.
func (bot *botClient) SetAlertConfig(alertConfig *protocol.AlertConfig) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	bot.alertConfigUnsafe = *alertConfig
}

// grpcClient returns the bot gRPC client.
func (bot *botClient) grpcClient() agentgrpc.Client {
	bot.mu.RLock()
	defer bot.mu.RUnlock()

	return bot.clientUnsafe
}

// setGrpcClient sets the bot config.
func (bot *botClient) setGrpcClient(client agentgrpc.Client) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	if bot.clientUnsafe != nil {
		bot.clientUnsafe.Close()
	}

	bot.clientUnsafe = client
}

// TxRequestCh returns the transaction request channel safely.
func (bot *botClient) TxRequestCh() chan<- *botreq.TxRequest {
	return bot.txRequests
}

// BlockRequestCh returns the block request channel safely.
func (bot *botClient) BlockRequestCh() chan<- *botreq.BlockRequest {
	return bot.blockRequests
}

// CombinationRequestCh returns the alert request channel safely.
func (bot *botClient) CombinationRequestCh() chan<- *botreq.CombinationRequest {
	return bot.combinationRequests
}

// Close implements io.Closer.
func (bot *botClient) Close() error {
	bot.closeOnce.Do(func() {
		bot.ctxCancel()
		client := bot.grpcClient()
		if client != nil {
			client.Close()
		}
		botConfig := bot.Config()
		log.WithField("bot", botConfig.ID).WithField("image", botConfig.Image).Info("detached")
		bot.lifecycleMetrics.ClientClose(botConfig)
		if bot.isCombinerBot() {
			bot.msgClient.Publish(messaging.SubjectAgentsAlertUnsubscribe, bot.CombinerBotSubscriptions())
			bot.lifecycleMetrics.ActionUnsubscribe(bot.CombinerBotSubscriptions())
		}
	})
	return nil
}

// Closed returns the closed channel.
func (bot *botClient) Closed() <-chan struct{} {
	return bot.ctx.Done()
}

// IsClosed tells if the bot is closed.
func (bot *botClient) IsClosed() bool {
	return isChanClosed(bot.ctx.Done())
}

// setInitialized sets the bot as initialized.
func (bot *botClient) setInitialized() {
	bot.initializedOnce.Do(
		func() {
			close(bot.initialized) // never close this anywhere else
		},
	)
}

// Initialized returns the initialized channel.
func (bot *botClient) Initialized() <-chan struct{} {
	return bot.initialized
}

// IsInitialized tells if the bot is initialized.
func (bot *botClient) IsInitialized() bool {
	return isChanClosed(bot.initialized)
}

func isChanClosed(ch <-chan struct{}) bool {
	select {
	case _, ok := <-ch:
		return !ok
	default:
		return false
	}
}

// Initialize initializes the bot.
func (bot *botClient) Initialize() {
	bot.initialize()
}

func (bot *botClient) initialize() {
	botConfig := bot.Config()

	logger := log.WithFields(log.Fields{
		"bot": botConfig.ID,
	})

	// publish start metric to track bot starts/restarts.
	bot.lifecycleMetrics.ClientDial(botConfig)

	botClient, err := bot.dialer.DialBot(botConfig)
	if err != nil {
		logger.WithError(err).Info("failed to dial bot")
		return
	}
	bot.setGrpcClient(botClient)
	bot.lifecycleMetrics.StatusAttached(botConfig)
	logger.Info("attached to bot")

	ctx, cancel := context.WithTimeout(bot.ctx, DefaultInitializeTimeout)
	defer cancel()

	// invoke initialize method of the bot
	initializeResponse, err := botClient.Initialize(ctx, &protocol.InitializeRequest{
		AgentId:   botConfig.ID,
		ProxyHost: config.DockerJSONRPCProxyContainerName,
	})

	// it is not mandatory to implement a initialize method, safe to skip
	if status.Code(err) == codes.Unimplemented {
		logger.WithError(err).Info("initialize() method not implemented in bot - safe to ignore")
		bot.initSuccess(botConfig)
		return
	}
	if err != nil {
		logger.WithError(err).Warn("bot initialization failed")
		bot.lifecycleMetrics.FailureInitialize(err, botConfig)
		_ = bot.Close()
		return
	}

	if initializeResponse != nil && initializeResponse.Status == protocol.ResponseStatus_ERROR {
		err := agentgrpc.Error(initializeResponse.Errors)
		logger.WithError(err).Warn("bot initialization returned an error response")
		bot.lifecycleMetrics.FailureInitializeResponse(err, botConfig)
		_ = bot.Close()
		return
	}

	if err := validateInitializeResponse(initializeResponse); err != nil {
		logger.WithError(err).Warn("bot initialization validation failed")
		bot.lifecycleMetrics.FailureInitializeValidate(err, botConfig)
		return
	}

	// Let services know about the latest subscriptions
	if initializeResponse != nil && initializeResponse.AlertConfig != nil {
		bot.SetAlertConfig(initializeResponse.AlertConfig)
		bot.msgClient.Publish(messaging.SubjectAgentsAlertSubscribe, bot.CombinerBotSubscriptions())
		bot.lifecycleMetrics.ActionSubscribe(bot.CombinerBotSubscriptions())
	}

	bot.initSuccess(botConfig)
	logger.Info("bot initialization succeeded")
}

func (bot *botClient) initSuccess(botConfig config.AgentConfig) {
	bot.setInitialized()
	bot.lifecycleMetrics.StatusInitialized(botConfig)
}

func validateInitializeResponse(response *protocol.InitializeResponse) error {
	if response == nil {
		return fmt.Errorf("initialize response can not be nil")
	}
	if response.AlertConfig == nil {
		return nil
	}

	for _, subscription := range response.AlertConfig.Subscriptions {
		if !utils.IsValidBotID(subscription.BotId) {
			return fmt.Errorf("invalid bot id: %s", subscription.BotId)
		}
	}

	return nil
}

// StartProcessing launches the goroutines to concurrently process incoming requests
// from request channels.
func (bot *botClient) StartProcessing() {
	go bot.processTransactions()
	go bot.processBlocks()
	go bot.processCombinationAlerts()
	go bot.processHealthChecks()
}

func processRequests[R any](
	ctx context.Context, reqCh <-chan *R, closedCh <-chan struct{}, logger *log.Entry,
	processFunc func(context.Context, *log.Entry, *R) bool,
) {
	for {
		select {
		case <-ctx.Done():
			logger.WithError(ctx.Err()).Info("bot context is done")
			return

		case request := <-reqCh:
			ctx, cancel := context.WithTimeout(ctx, RequestTimeout)
			exit := processFunc(ctx, logger, request)
			cancel()
			if exit {
				return
			}
		}
	}
}

func (bot *botClient) processTransactions() {
	lg := log.WithFields(
		log.Fields{
			"bot":       bot.Config().ID,
			"component": "bot-client",
			"evaluate":  "transaction",
		},
	)

	<-bot.Initialized()

	processRequests(bot.ctx, bot.txRequests, bot.Closed(), lg, bot.processTransaction)
}

func (bot *botClient) processBlocks() {
	lg := log.WithFields(
		log.Fields{
			"bot":       bot.Config().ID,
			"component": "bot-client",
			"evaluate":  "block",
		},
	)

	<-bot.Initialized()

	processRequests(bot.ctx, bot.blockRequests, bot.Closed(), lg, bot.processBlock)
}

func (bot *botClient) processHealthChecks() {
	lg := log.WithFields(
		log.Fields{
			"bot":       bot.Config().ID,
			"component": "bot-client",
			"evaluate":  "health-check",
		},
	)

	<-bot.Initialized()

	t := time.NewTicker(DefaultHealthCheckInterval)

	for {
		select {
		case <-bot.ctx.Done():
			return
		case <-t.C:
			exit := bot.doHealthCheck(bot.ctx, lg)
			if exit {
				return
			}
		}
	}
}

func (bot *botClient) processCombinationAlerts() {
	lg := log.WithFields(
		log.Fields{
			"bot":       bot.Config().ID,
			"component": "bot-client",
			"evaluate":  "combination",
		},
	)

	<-bot.Initialized()

	processRequests(bot.ctx, bot.combinationRequests, bot.Closed(), lg, bot.processCombinationAlert)
}

func (bot *botClient) processTransaction(ctx context.Context, lg *log.Entry, request *botreq.TxRequest) (exit bool) {
	botConfig := bot.Config()
	botClient := bot.grpcClient()

	if bot.IsClosed() {
		return true
	}

	startTime := time.Now()

	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateTxResponse)

	requestTime := time.Now().UTC()
	err := botClient.Invoke(ctx, agentgrpc.MethodEvaluateTx, request.Original, resp)
	responseTime := time.Now().UTC()

	if err == nil {
		// truncate findings
		if len(resp.Findings) > MaxFindings {
			dropped := len(resp.Findings) - MaxFindings
			droppedMetric := metrics.CreateAgentMetric(botConfig, metrics.MetricFindingsDropped, float64(dropped), "")
			bot.msgClient.PublishProto(
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
		resp.Metadata["imageHash"] = botConfig.ImageHash()

		ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
		ts.BotRequest = requestTime
		ts.BotResponse = responseTime

		bot.resultChannels.Tx <- &botreq.TxResult{
			AgentConfig: botConfig,
			Request:     request.Original,
			Response:    resp,
			Timestamps:  ts,
		}
		lg.WithField("duration", time.Since(startTime)).Debugf("sent results")

		return false
	}

	if status.Code(err) == codes.Unimplemented {
		return false
	}

	lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking bot")
	bot.lifecycleMetrics.BotError("tx.invoke", err, botConfig)

	if bot.errCounter.TooManyErrs(err) {
		lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down bot")
		_ = bot.Close()
		bot.lifecycleMetrics.FailureTooManyErrs(err, botConfig)
		return true
	}

	return false
}

func (bot *botClient) processBlock(ctx context.Context, lg *log.Entry, request *botreq.BlockRequest) (exit bool) {
	botConfig := bot.Config()
	botClient := bot.grpcClient()

	if bot.IsClosed() {
		return true
	}

	startTime := time.Now()

	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateBlockResponse)
	requestTime := time.Now().UTC()
	err := botClient.Invoke(ctx, agentgrpc.MethodEvaluateBlock, request.Original, resp)
	responseTime := time.Now().UTC()

	if err == nil {
		// truncate findings
		if len(resp.Findings) > MaxFindings {
			dropped := len(resp.Findings) - MaxFindings
			droppedMetric := metrics.CreateAgentMetric(
				botConfig, metrics.MetricFindingsDropped, float64(dropped), "",
			)
			bot.msgClient.PublishProto(
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
		resp.Metadata["imageHash"] = botConfig.ImageHash()

		ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
		ts.BotRequest = requestTime
		ts.BotResponse = responseTime

		bot.resultChannels.Block <- &botreq.BlockResult{
			AgentConfig: botConfig,
			Request:     request.Original,
			Response:    resp,
			Timestamps:  ts,
		}
		lg.WithField("duration", time.Since(startTime)).Debugf("sent results")

		return false
	}

	if status.Code(err) == codes.Unimplemented {
		return false
	}

	lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking bot")
	bot.lifecycleMetrics.BotError("block.invoke", err, botConfig)

	if bot.errCounter.TooManyErrs(err) {
		lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down bot")
		_ = bot.Close()
		bot.lifecycleMetrics.FailureTooManyErrs(err, botConfig)
		return true
	}

	return false
}

func (bot *botClient) processCombinationAlert(ctx context.Context, lg *log.Entry, request *botreq.CombinationRequest) bool {
	botConfig := bot.Config()
	botClient := bot.grpcClient()

	if bot.IsClosed() {
		return true
	}

	startTime := time.Now()

	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")
	resp := new(protocol.EvaluateAlertResponse)
	requestTime := time.Now().UTC()
	err := botClient.Invoke(ctx, agentgrpc.MethodEvaluateAlert, request.Original, resp)
	responseTime := time.Now().UTC()

	if err != nil {
		if status.Code(err) != codes.Unimplemented {
			lg.WithField("duration", time.Since(startTime)).WithError(err).Error("error invoking bot")
			bot.lifecycleMetrics.BotError("combiner.invoke", err, botConfig)
		}


		if bot.errCounter.TooManyErrs(err) {
			lg.WithField("duration", time.Since(startTime)).Error("too many errors - shutting down bot")
			_ = bot.Close()
			bot.lifecycleMetrics.FailureTooManyErrs(err, botConfig)
			return true
		}
	}

	// validate response
	if vErr := validateEvaluateAlertResponse(resp); vErr != nil {
		lg.WithField(
			"request", request.Original.RequestId,
		).WithError(vErr).Error("evaluate combination response validation failed")
		bot.lifecycleMetrics.BotError("validate.evaluate.alert.response", vErr, botConfig)

		return false
	}

	// truncate findings
	if len(resp.Findings) > MaxFindings {
		dropped := len(resp.Findings) - MaxFindings
		droppedMetric := metrics.CreateAgentMetric(botConfig, metrics.MetricFindingsDropped, float64(dropped), "")
		bot.msgClient.PublishProto(
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

	resp.Metadata["imageHash"] = botConfig.ImageHash()

	ts := domain.TrackingTimestampsFromMessage(request.Original.Event.Timestamps)
	ts.BotRequest = requestTime
	ts.BotResponse = responseTime

	bot.resultChannels.CombinationAlert <- &botreq.CombinationAlertResult{
		AgentConfig: botConfig,
		Request:     request.Original,
		Response:    resp,
		Timestamps:  ts,
	}

	lg.WithField("duration", time.Since(startTime)).Debugf("sent results")
	return false
}

func (bot *botClient) doHealthCheck(ctx context.Context, lg *log.Entry) bool {
	botConfig := bot.Config()
	botClient := bot.grpcClient()

	if bot.IsClosed() {
		return true
	}

	startTime := time.Now()

	lg.WithField("duration", time.Since(startTime)).Debugf("sending request")

	bot.lifecycleMetrics.HealthCheckAttempt(botConfig)

	err := botClient.DoHealthCheck(ctx)
	if err != nil {
		bot.lifecycleMetrics.HealthCheckError(err, botConfig)
	} else {
		bot.lifecycleMetrics.HealthCheckSuccess(botConfig)
	}

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

// ShouldProcessBlock tells if the bot should process block.
func (bot *botClient) ShouldProcessBlock(blockNumberHex string) bool {
	botConfig := bot.Config()

	blockNumber, _ := hexutil.DecodeUint64(blockNumberHex)
	var isAtLeastStartBlock bool
	if botConfig.StartBlock != nil {
		isAtLeastStartBlock = blockNumber >= *botConfig.StartBlock
	} else {
		isAtLeastStartBlock = true
	}

	var isAtMostStopBlock bool
	if botConfig.StopBlock != nil {
		isAtMostStopBlock = blockNumber <= *botConfig.StopBlock
	} else {
		isAtMostStopBlock = true
	}

	var isOnThisShard bool
	// if sharded, block % shards must be equal to shard id
	if botConfig.IsSharded() {
		isOnThisShard = uint(blockNumber)%botConfig.ShardConfig.Shards == botConfig.ShardConfig.ShardID
	} else {
		isOnThisShard = true
	}

	return isAtLeastStartBlock && isAtMostStopBlock && isOnThisShard
}

func (bot *botClient) ShouldProcessAlert(event *protocol.AlertEvent) bool {
	if !bot.isCombinerBot() {
		return false
	}

	botConfig := bot.Config()

	// handle sharding
	alertCreatedAt, err := time.Parse(time.RFC3339Nano, event.Alert.CreatedAt)
	if err != nil {
		log.WithFields(
			log.Fields{
				"alertHash": event.Alert.Hash,
				"createdAt": event.Alert.CreatedAt,
				"botId":     botConfig.ID,
			},
		).Warn("failed to parse created at for sharding calculation")

		return false
	}

	var isOnThisShard bool
	if botConfig.IsSharded() {
		isOnThisShard = uint(alertCreatedAt.Unix())%botConfig.ShardConfig.Shards == botConfig.ShardConfig.ShardID
	} else {
		isOnThisShard = true
	}

	return isOnThisShard
}
