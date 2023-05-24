package botio

import (
	"context"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/services/components/botio/botreq"
	"github.com/forta-network/forta-node/services/components/metrics"
	log "github.com/sirupsen/logrus"
)

// Sender sends requests to all bots and outputs bot responses.
type Sender interface {
	SendEvaluateTxRequest(req *protocol.EvaluateTxRequest)
	SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest)
	SendEvaluateAlertRequest(req *protocol.EvaluateAlertRequest)
	health.Reporter
}

// BotPool knows the latest bot clients.
type BotPool interface {
	WaitForAll()
	GetCurrentBotClients() []BotClient
}

type requestSender struct {
	ctx context.Context

	botPool   BotPool
	msgClient clients.MessageClient
}

// NewSender creates a new requestSender.
func NewSender(ctx context.Context, msgClient clients.MessageClient, botPool BotPool) Sender {
	return &requestSender{
		ctx:       ctx,
		botPool:   botPool,
		msgClient: msgClient,
	}
}

// Health implements health.Reporter interface.
func (rs *requestSender) Health() health.Reports {
	bots := rs.botPool.GetCurrentBotClients()

	botCount := len(bots)
	var fullCount int
	for _, bot := range bots {
		if bot.TxBufferIsFull() {
			fullCount++
		}
	}
	status := health.StatusOK
	if botCount == 0 {
		status = health.StatusFailing
	}
	return health.Reports{
		&health.Report{
			Name:    "agents.total",
			Status:  status,
			Details: strconv.Itoa(botCount),
		},
		&health.Report{
			Name:    "agents.lagging",
			Status:  health.StatusInfo,
			Details: strconv.Itoa(fullCount),
		},
	}
}

// Name implements health.Reporter interface.
func (rs *requestSender) Name() string {
	return "sender"
}

// SendEvaluateTxRequest sends the request to all of the active bots which
// should be processing the block.
func (rs *requestSender) SendEvaluateTxRequest(req *protocol.EvaluateTxRequest) {
	startTime := time.Now()
	lg := log.WithFields(log.Fields{
		"tx":        req.Event.Transaction.Hash,
		"component": "pool",
	})
	lg.Debug("SendEvaluateTxRequest")

	rs.botPool.WaitForAll()

	bots := rs.botPool.GetCurrentBotClients()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}
	var metricsList []*protocol.AgentMetric
	for _, bot := range bots {
		if !bot.ShouldProcessBlock(req.Event.Block.BlockNumber) {
			continue
		}
		lg.WithFields(log.Fields{
			"bot":      bot.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending tx request to evalTxCh")

		// unblock req send and discard agent if agent is closed

		select {
		case <-bot.Closed():
			lg.WithField("bot", bot.Config().ID).Debug("bot is closed - skipping")
		case bot.TxRequestCh() <- &botreq.TxRequest{
			Original: req,
			Encoded:  encoded,
		}:
		default: // do not try to send if the buffer is full
			lg.WithField("bot", bot.Config().ID).Debug("agent tx request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(bot.Config().ID, metrics.MetricTxDrop, 1))
		}
		lg.WithFields(log.Fields{
			"bot":      bot.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sent tx request to evalTxCh")
	}
	metrics.SendAgentMetrics(rs.msgClient, metricsList)

	lg.WithFields(log.Fields{
		"duration": time.Since(startTime),
	}).Debug("Finished SendEvaluateTxRequest")
}

// SendEvaluateBlockRequest sends the request to all of the active bots which
// should be processing the block.
func (rs *requestSender) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	startTime := time.Now()
	lg := log.WithFields(log.Fields{
		"block":     req.Event.BlockNumber,
		"component": "pool",
	})
	lg.Debug("SendEvaluateBlockRequest")

	rs.botPool.WaitForAll()

	bots := rs.botPool.GetCurrentBotClients()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}

	var metricsList []*protocol.AgentMetric
	for _, bot := range bots {
		if !bot.ShouldProcessBlock(req.Event.BlockNumber) {
			continue
		}

		lg.WithFields(log.Fields{
			"bot":      bot.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending block request to evalBlockCh")

		// unblock req send if agent is closed
		select {
		case <-bot.Closed():
			lg.WithField("bot", bot.Config().ID).Debug("bot is closed - skipping")
		case bot.BlockRequestCh() <- &botreq.BlockRequest{
			Original: req,
			Encoded:  encoded,
		}:
		default: // do not try to send if the buffer is full
			lg.WithField("bot", bot.Config().ID).Warn("agent block request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(bot.Config().ID, metrics.MetricBlockDrop, 1))
		}
		lg.WithFields(
			log.Fields{
				"bot":      bot.Config().ID,
				"duration": time.Since(startTime),
			},
		).Debug("sent tx request to evalBlockCh")
	}

	blockNumber, _ := hexutil.DecodeUint64(req.Event.BlockNumber)
	rs.msgClient.Publish(messaging.SubjectScannerBlock, &messaging.ScannerPayload{
		LatestBlockInput: blockNumber,
	})

	metrics.SendAgentMetrics(rs.msgClient, metricsList)
	lg.WithFields(log.Fields{
		"duration": time.Since(startTime),
	}).Debug("Finished SendEvaluateBlockRequest")
}

// SendEvaluateAlertRequest sends the request to all the active bots which
// should be processing the alert.
func (rs *requestSender) SendEvaluateAlertRequest(req *protocol.EvaluateAlertRequest) {
	startTime := time.Now()
	lg := log.WithFields(
		log.Fields{
			"target":    req.TargetBotId,
			"alert":     req.Event.Alert.Hash,
			"component": "pool",
		},
	)
	lg.Debug("SendEvaluateAlertRequest")

	if req.Event.Alert == nil || req.Event.Alert.Source == nil || req.Event.Alert.Source.Bot == nil {
		lg.Warn("bad request")
		return
	}

	rs.botPool.WaitForAll()

	bots := rs.botPool.GetCurrentBotClients()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}

	var metricsList []*protocol.AgentMetric

	var target BotClient

	// find target bot for the event
	for _, bot := range bots {
		if bot.Config().ID != req.TargetBotId {
			continue
		}

		if !bot.IsInitialized() {
			continue
		}
		target = bot
		break
	}

	// return if can't find the target bot, or it's not ready yet
	if target == nil {
		lg.Warn("failed to find subscriber")
		return
	}

	// filter out bad events
	if !target.ShouldProcessAlert(req.Event) {
		return
	}

	lg.WithFields(
		log.Fields{
			"bot":      target.Config().ID,
			"duration": time.Since(startTime),
		},
	).Debug("sending alert request to evalAlertCh")

	// unblock req send if agent is closed
	select {
	case <-target.Closed():
		lg.WithField("bot", target.Config().ID).Debug("bot is closed - skipping")
	case target.CombinationRequestCh() <- &botreq.CombinationRequest{
		Original: req,
		Encoded:  encoded,
	}:
	default: // do not try to send if the buffer is full
		lg.WithField("bot", target.Config().ID).Warn("agent alert request buffer is full - skipping")
		metricsList = append(metricsList, metrics.CreateAgentMetric(target.Config().ID, metrics.MetricCombinerDrop, 1))
	}

	lg.WithFields(
		log.Fields{
			"bot":      target.Config().ID,
			"duration": time.Since(startTime),
		},
	).Debug("sent alert request to evalAlertCh")

	rs.msgClient.Publish(messaging.SubjectScannerAlert, &messaging.ScannerPayload{})
	metrics.SendAgentMetrics(rs.msgClient, metricsList)
	lg.WithFields(
		log.Fields{
			"duration": time.Since(startTime),
		},
	).Debug("Finished SendEvaluateAlertRequest")
}
