package agentpool

import (
	"sync"
	"time"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/protocol"
	"github.com/forta-network/forta-node/services/scanner"
	"github.com/forta-network/forta-node/services/scanner/agentpool/poolagent"
	log "github.com/sirupsen/logrus"
)

// Constants
const (
	DefaultBufferSize = 100
)

// AgentPool maintains the pool of agents that the scanner should
// interact with.
type AgentPool struct {
	agents       []*poolagent.Agent
	txResults    chan *scanner.TxResult
	blockResults chan *scanner.BlockResult
	msgClient    clients.MessageClient
	dialer       func(config.AgentConfig) clients.AgentClient
	mu           sync.RWMutex
}

// NewAgentPool creates a new agent pool.
func NewAgentPool(msgClient clients.MessageClient) *AgentPool {
	agentPool := &AgentPool{
		txResults:    make(chan *scanner.TxResult, DefaultBufferSize),
		blockResults: make(chan *scanner.BlockResult, DefaultBufferSize),
		msgClient:    msgClient,
		dialer: func(ac config.AgentConfig) clients.AgentClient {
			client := agentgrpc.NewClient()
			client.MustDial(ac)
			return client
		},
	}
	agentPool.registerMessageHandlers()
	go agentPool.logAgentChanBuffersLoop()
	return agentPool
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
	log.WithField("tx", req.Event.Transaction.Hash).Debug("SendEvaluateTxRequest")

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.Block.BlockNumber) {
			continue
		}
		log.WithField("agent", agent.Config().ID).Debug("sending tx request to evalBlockCh")

		// unblock req send and discard agent if agent is closed
		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
		case agent.TxRequestCh() <- req:
		default: // do not try to send if the buffer is full
		}
	}
	log.WithField("tx", req.Event.Transaction.Hash).Debug("Finished SendEvaluateTxRequest")
}

// TxResults returns the receive-only tx results channel.
func (ap *AgentPool) TxResults() <-chan *scanner.TxResult {
	return ap.txResults
}

// SendEvaluateBlockRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	log.WithField("block", req.Event.BlockNumber).Debug("SendEvaluateBlockRequest")

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

	for _, agent := range agents {
		if !agent.IsReady() || !agent.ShouldProcessBlock(req.Event.BlockNumber) {
			continue
		}
		log.WithField("agent", agent.Config().ID).Debug("sending block request to evalBlockCh")

		// unblock req send if agent is closed
		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
		case agent.BlockRequestCh() <- req:
		default: // do not try to send if the buffer is full
		}
	}
	log.WithField("block", req.Event.BlockNumber).Debug("Finished SendEvaluateBlockRequest")
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
	for _, agentCfg := range payload {
		for _, agent := range ap.agents {
			if agent.Config().ContainerName() == agentCfg.ContainerName() {
				agent.SetClient(ap.dialer(agent.Config()))
				agent.SetReady()
				agent.StartProcessing()
			}
		}
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
				log.WithField("agent", agent.Config().ID).Debug("stopping")
				agent.Close()
				stopped = true
				break
			}
		}
		if !stopped {
			log.WithField("agent", agent.Config().ID).Debug("not stopped")
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
