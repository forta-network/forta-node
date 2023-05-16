package agentpool

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/metrics"
	"github.com/forta-network/forta-node/services/scanner"
	"github.com/forta-network/forta-node/services/scanner/agentpool/poolagent"
	log "github.com/sirupsen/logrus"
)

type Dispatcher struct {
	agents                  []*poolagent.Agent
	txResults               chan *scanner.TxResult
	blockResults            chan *scanner.BlockResult
	combinationAlertResults chan *scanner.CombinationAlertResult
	mu                      sync.RWMutex
	msgClient               clients.MessageClient
}

func (d *Dispatcher) SetAgents(agents []*poolagent.Agent) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.agents = agents
}

// SendEvaluateTxRequest sends the request to all of the active agents which
// should be processing the block.
func (d *Dispatcher) SendEvaluateTxRequest(req *protocol.EvaluateTxRequest) {
	startTime := time.Now()
	lg := log.WithFields(
		log.Fields{
			"tx":        req.Event.Transaction.Hash,
			"component": "pool",
		},
	)
	lg.Debug("SendEvaluateTxRequest")

	d.mu.RLock()
	agents := d.agents
	d.mu.RUnlock()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}
	var metricsList []*protocol.AgentMetric
	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.Block.BlockNumber) {
			continue
		}
		lg.WithFields(
			log.Fields{
				"agent":    agent.Config().ID,
				"duration": time.Since(startTime),
			},
		).Debug("sending tx request to evalTxCh")

		// unblock req send and discard agent if agent is closed

		select {
		// TODO: closed agent shouldn't exist in dispatcher
		case <-agent.Closed():
		// 	d.discardAgent(agent)
		case agent.TxRequestCh() <- &poolagent.TxRequest{
			Original: req,
			Encoded:  encoded,
		}:
		default: // do not try to send if the buffer is full
			lg.WithField("agent", agent.Config().ID).Debug("agent tx request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(agent.Config().ID, metrics.MetricTxDrop, 1))
		}
		lg.WithFields(
			log.Fields{
				"agent":    agent.Config().ID,
				"duration": time.Since(startTime),
			},
		).Debug("sent tx request to evalTxCh")
	}
	metrics.SendAgentMetrics(d.msgClient, metricsList)

	lg.WithFields(
		log.Fields{
			"duration": time.Since(startTime),
		},
	).Debug("Finished SendEvaluateTxRequest")
}

// TxResults returns the receive-only tx results channel.
func (d *Dispatcher) TxResults() <-chan *scanner.TxResult {
	return d.txResults
}

// SendEvaluateBlockRequest sends the request to all of the active agents which
// should be processing the block.
func (d *Dispatcher) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	startTime := time.Now()
	lg := log.WithFields(
		log.Fields{
			"block":     req.Event.BlockNumber,
			"component": "pool",
		},
	)
	lg.Debug("SendEvaluateBlockRequest")

	d.mu.RLock()
	agents := d.agents
	d.mu.RUnlock()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}

	var metricsList []*protocol.AgentMetric
	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.BlockNumber) {
			continue
		}

		lg.WithFields(
			log.Fields{
				"agent":    agent.Config().ID,
				"duration": time.Since(startTime),
			},
		).Debug("sending block request to evalBlockCh")

		// unblock req send if agent is closed
		select {
		// TODO: closed agent shouldn't exist in dispatcher
		case <-agent.Closed():
		// 	d.discardAgent(agent)
		case agent.BlockRequestCh() <- &poolagent.BlockRequest{
			Original: req,
			Encoded:  encoded,
		}:
		default: // do not try to send if the buffer is full
			lg.WithField("agent", agent.Config().ID).Warn("agent block request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(agent.Config().ID, metrics.MetricBlockDrop, 1))
		}
		lg.WithFields(
			log.Fields{
				"agent":    agent.Config().ID,
				"duration": time.Since(startTime),
			},
		).Debug("sent tx request to evalBlockCh")
	}

	blockNumber, _ := hexutil.DecodeUint64(req.Event.BlockNumber)
	d.msgClient.Publish(
		messaging.SubjectScannerBlock, &messaging.ScannerPayload{
			LatestBlockInput: blockNumber,
		},
	)

	metrics.SendAgentMetrics(d.msgClient, metricsList)
	lg.WithFields(
		log.Fields{
			"duration": time.Since(startTime),
		},
	).Debug("Finished SendEvaluateBlockRequest")
}

// SendEvaluateAlertRequest sends the request to all the active agents which
// should be processing the alert.
func (d *Dispatcher) SendEvaluateAlertRequest(req *protocol.EvaluateAlertRequest) {
	startTime := time.Now()
	lg := log.WithFields(
		log.Fields{
			"component": "pool",
			"target":    req.TargetBotId,
		},
	)
	lg.Debug("SendEvaluateAlertRequest")

	if req.Event.Alert == nil || req.Event.Alert.Source == nil || req.Event.Alert.Source.Bot == nil {
		lg.Warn("bad request")
		return
	}

	d.mu.RLock()
	agents := d.agents
	d.mu.RUnlock()

	encoded, err := agentgrpc.EncodeMessage(req)
	if err != nil {
		lg.WithError(err).Error("failed to encode message")
		return
	}

	var metricsList []*protocol.AgentMetric

	var target *poolagent.Agent

	// find target bot for the event
	for _, agent := range agents {
		if agent.Config().ID != req.TargetBotId {
			continue
		}

		if !agent.IsReady() {
			continue
		}
		target = agent
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
			"agent":    target.Config().ID,
			"duration": time.Since(startTime),
		},
	).Debug("sending alert request to evalAlertCh")

	// unblock req send if agent is closed
	select {
	// TODO: closed agent shouldn't exist in dispatcher
	case <-target.Closed():
	// 	d.discardAgent(agent)
	case target.CombinationRequestCh() <- &poolagent.CombinationRequest{
		Original: req,
		Encoded:  encoded,
	}:
	default: // do not try to send if the buffer is full
		lg.WithField("agent", target.Config().ID).Warn("agent alert request buffer is full - skipping")
		metricsList = append(metricsList, metrics.CreateAgentMetric(target.Config().ID, metrics.MetricCombinerDrop, 1))
	}

	lg.WithFields(
		log.Fields{
			"agent":    target.Config().ID,
			"duration": time.Since(startTime),
		},
	).Debug("sent alert request to evalAlertCh")

	d.msgClient.Publish(messaging.SubjectScannerAlert, &messaging.ScannerPayload{})
	metrics.SendAgentMetrics(d.msgClient, metricsList)
	lg.WithFields(
		log.Fields{
			"duration": time.Since(startTime),
		},
	).Debug("Finished SendEvaluateAlertRequest")
}

// CombinationAlertResults returns the receive-only alert results channel.
func (d *Dispatcher) CombinationAlertResults() <-chan *scanner.CombinationAlertResult {
	return d.combinationAlertResults
}

// BlockResults returns the receive-only tx results channel.
func (d *Dispatcher) BlockResults() <-chan *scanner.BlockResult {
	return d.blockResults
}
