package agentpool

import (
	"strconv"
	"sync"
	"time"

	"github.com/forta-protocol/forta-node/metrics"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/agentgrpc"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/protocol"
	"github.com/forta-protocol/forta-node/services/scanner"
	"github.com/forta-protocol/forta-node/services/scanner/agentpool/poolagent"
	log "github.com/sirupsen/logrus"
)

// Constants
const (
	DefaultBufferSize = 100 * poolagent.DefaultBufferSize // i.e. assuming 100 agents
)

// AgentPool maintains the pool of agents that the scanner should
// interact with.
type AgentPool struct {
	agents       []*poolagent.Agent
	txResults    chan *scanner.TxResult
	blockResults chan *scanner.BlockResult
	msgClient    clients.MessageClient
	dialer       func(config.AgentConfig) (clients.AgentClient, error)
	mu           sync.RWMutex
}

// NewAgentPool creates a new agent pool.
func NewAgentPool(cfg config.ScannerConfig, msgClient clients.MessageClient) *AgentPool {
	agentPool := &AgentPool{
		txResults:    make(chan *scanner.TxResult, DefaultBufferSize),
		blockResults: make(chan *scanner.BlockResult, DefaultBufferSize),
		msgClient:    msgClient,
		dialer: func(ac config.AgentConfig) (clients.AgentClient, error) {
			client := agentgrpc.NewClient()
			if err := client.Dial(ac); err != nil {
				return nil, err
			}
			return client, nil
		},
	}

	agentPool.registerMessageHandlers()
	go agentPool.logAgentChanBuffersLoop()
	return agentPool
}

// Health implements health.Reporter interface.
func (ap *AgentPool) Health() health.Reports {
	ap.mu.RLock()
	defer ap.mu.RUnlock()

	agentCount := len(ap.agents)
	var fullCount int
	for _, agent := range ap.agents {
		if agent.TxBufferIsFull() {
			fullCount++
		}
	}
	status := health.StatusOK
	if agentCount == 0 {
		status = health.StatusFailing
	}
	return health.Reports{
		&health.Report{
			Name:    "agents.total",
			Status:  status,
			Details: strconv.Itoa(agentCount),
		},
		&health.Report{
			Name:    "agents.lagging",
			Status:  health.StatusInfo,
			Details: strconv.Itoa(fullCount),
		},
	}
}

// Name implements health.Reporter interface.
func (ap *AgentPool) Name() string {
	return "agent-pool"
}

// discardAgent removes the agent from the list which eventually causes the
// request channels to be deallocated.
func (ap *AgentPool) discardAgent(discarded *poolagent.Agent) {
	ap.mu.Lock()
	var newAgents []*poolagent.Agent
	for _, agent := range ap.agents {
		if agent != discarded {
			newAgents = append(newAgents, agent)
		} else {
			log.WithField("agent", agent.Config().ContainerName()).Info("discarded")
		}
	}
	ap.agents = newAgents
	ap.mu.Unlock()
}

// SendEvaluateTxRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateTxRequest(req *protocol.EvaluateTxRequest) {
	startTime := time.Now()
	lg := log.WithFields(log.Fields{
		"tx":        req.Event.Transaction.Hash,
		"component": "pool",
	})
	lg.Debug("SendEvaluateTxRequest")

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

	var metricsList []*protocol.AgentMetric
	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.Block.BlockNumber) {
			continue
		}
		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending tx request to evalTxCh")

		// unblock req send and discard agent if agent is closed
		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
		case agent.TxRequestCh() <- req:
		default: // do not try to send if the buffer is full
			lg.WithField("agent", agent.Config().ID).Warn("agent tx request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(agent.Config().ID, metrics.MetricTxDrop, 1))
		}
		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sent tx request to evalTxCh")
	}
	metrics.SendAgentMetrics(ap.msgClient, metricsList)
	lg.WithFields(log.Fields{
		"duration": time.Since(startTime),
	}).Debug("Finished SendEvaluateTxRequest")
}

// TxResults returns the receive-only tx results channel.
func (ap *AgentPool) TxResults() <-chan *scanner.TxResult {
	return ap.txResults
}

// SendEvaluateBlockRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	startTime := time.Now()
	lg := log.WithFields(log.Fields{
		"block":     req.Event.BlockNumber,
		"component": "pool",
	})
	lg.Debug("SendEvaluateBlockRequest")
	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

	var metricsList []*protocol.AgentMetric
	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.BlockNumber) {
			continue
		}

		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending block request to evalBlockCh")

		// unblock req send if agent is closed
		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
		case agent.BlockRequestCh() <- req:
		default: // do not try to send if the buffer is full
			lg.WithField("agent", agent.Config().ID).Warn("agent block request buffer is full - skipping")
			metricsList = append(metricsList, metrics.CreateAgentMetric(agent.Config().ID, metrics.MetricBlockDrop, 1))
		}
		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sent tx request to evalBlockCh")
	}
	metrics.SendAgentMetrics(ap.msgClient, metricsList)
	lg.WithFields(log.Fields{
		"duration": time.Since(startTime),
	}).Debug("Finished SendEvaluateBlockRequest")
}

func (ap *AgentPool) logAgentChanBuffersLoop() {
	ticker := time.NewTicker(time.Second * 30)
	for range ticker.C {
		ap.logAgentStatuses()
	}
}

func (ap *AgentPool) logAgentStatuses() {
	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

	for _, agent := range agents {
		agent.LogStatus()
	}
}

// BlockResults returns the receive-only tx results channel.
func (ap *AgentPool) BlockResults() <-chan *scanner.BlockResult {
	return ap.blockResults
}

func (ap *AgentPool) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	log.Debug("handleAgentVersionsUpdate")
	latestVersions := payload

	// The agents list which we completely replace with the old ones.
	var newAgents []*poolagent.Agent

	// Find the missing agents in the pool, add them to the new agents list
	// and send a "run" message.
	var agentsToRun []config.AgentConfig
	for _, agentCfg := range latestVersions {
		var found bool
		for _, agent := range ap.agents {
			found = found || (agent.Config().ContainerName() == agentCfg.ContainerName())
		}
		if !found {
			newAgents = append(newAgents, poolagent.New(agentCfg, ap.msgClient, ap.txResults, ap.blockResults))
			agentsToRun = append(agentsToRun, agentCfg)
			log.WithField("agent", agentCfg.ID).Info("will trigger start")
		}
	}

	// Find the missing agents in the latest versions and send a "stop" message.
	// Otherwise, add to the new agents list so we keep on running.
	var agentsToStop []config.AgentConfig
	for _, agent := range ap.agents {
		var found bool
		var agentCfg config.AgentConfig
		for _, agentCfg = range latestVersions {
			found = found || (agent.Config().ContainerName() == agentCfg.ContainerName())
			if found {
				break
			}
		}
		if !found {
			agent.Close()
			agentsToStop = append(agentsToStop, agent.Config())
			log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Info("will trigger stop")
		} else {
			newAgents = append(newAgents, agent)
		}
	}

	ap.agents = newAgents
	if len(agentsToRun) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsActionRun, agentsToRun)
	}
	if len(agentsToStop) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}
	return nil
}

func (ap *AgentPool) handleStatusRunning(payload messaging.AgentPayload) error {
	log.Debug("handleStatusRunning")
	// If an agent was added before and just started to run, we should mark as ready.
	var agentsToStop []config.AgentConfig
	var agentsReady []config.AgentConfig

	for _, agentCfg := range payload {
		for _, agent := range ap.agents {
			if agent.Config().ContainerName() == agentCfg.ContainerName() {
				c, err := ap.dialer(agent.Config())
				if err != nil {
					log.WithField("agent", agent.Config().ID).WithError(err).Error("handleStatusRunning: error while dialing")
					agentsToStop = append(agentsToStop, agent.Config())
					continue
				}
				agent.SetClient(c)
				agent.SetReady()
				agent.StartProcessing()
				log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Info("attached")
				agentsReady = append(agentsReady, agent.Config())
			}
		}
	}
	if len(agentsReady) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsStatusAttached, agentsReady)
	}
	if len(agentsToStop) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}
	return nil
}

func (ap *AgentPool) handleStatusStopped(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	log.Debug("handleStatusStopped")
	var newAgents []*poolagent.Agent
	for _, agent := range ap.agents {
		var stopped bool
		for _, agentCfg := range payload {
			if agent.Config().ContainerName() == agentCfg.ContainerName() {
				agent.Close()
				log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Info("detached")
				stopped = true
				break
			}
		}
		if !stopped {
			log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Debug("not stopped")
			newAgents = append(newAgents, agent)
		}
	}
	ap.agents = newAgents
	return nil
}

func (ap *AgentPool) registerMessageHandlers() {
	ap.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(ap.handleAgentVersionsUpdate))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusRunning, messaging.AgentsHandler(ap.handleStatusRunning))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusStopped, messaging.AgentsHandler(ap.handleStatusStopped))
}
