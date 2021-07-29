package agentpool

import (
	"sync"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/agentgrpc"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
	"OpenZeppelin/fortify-node/services/scanner"
)

// AgentPool maintains the pool of agents that the scanner should
// interact with.
type AgentPool struct {
	agents       []*Agent
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
	return agentPool
}

// SendEvaluateTxRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateTxRequest(req *protocol.EvaluateTxRequest) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	agents := ap.agents
	for _, agent := range agents {
		if !agent.ready || !agent.shouldProcessBlock(req.Event.Block.BlockNumber) {
			continue
		}
		agent.evalTxCh <- req
	}
}

// TxResults returns the receive-only tx results channel.
func (ap *AgentPool) TxResults() <-chan *scanner.TxResult {
	return ap.txResults
}

// SendEvaluateBlockRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	agents := ap.agents
	for _, agent := range agents {
		if !agent.ready || !agent.shouldProcessBlock(req.Event.BlockNumber) {
			continue
		}
		agent.evalBlockCh <- req
	}
}

// BlockResults returns the receive-only tx results channel.
func (ap *AgentPool) BlockResults() <-chan *scanner.BlockResult {
	return ap.blockResults
}

func (ap *AgentPool) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	latestVersions := payload

	// The agents list which we completely replace with the old ones.
	var newAgents []*Agent

	// Find the missing agents in the pool, add them to the new agents list
	// and send a "run" message.
	var agentsToRun []config.AgentConfig
	for _, agentCfg := range latestVersions {
		var found bool
		for _, agent := range ap.agents {
			found = found || (agent.config.ContainerName() == agentCfg.ContainerName())
		}
		if !found {
			newAgents = append(newAgents, NewAgent(agentCfg, ap.txResults, ap.blockResults))
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
			found = found || (agent.config.ContainerName() == agentCfg.ContainerName())
			if found {
				break
			}
		}
		if !found {
			agent.Close()
			agent.ready = false
			agentsToStop = append(agentsToStop, agent.config)
		} else {
			newAgents = append(newAgents, agent)
		}
	}

	ap.agents = newAgents
	ap.manageReadinessUnsafe()
	if len(agentsToRun) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsActionRun, agentsToRun)
	}
	if len(agentsToStop) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}
	return nil
}

func (ap *AgentPool) handleStatusRunning(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	// If an agent was added before and just started to run, we should mark as ready
	// and start the processing goroutines.
	for _, agentCfg := range payload {
		for _, agent := range ap.agents {
			if agent.config.ContainerName() == agentCfg.ContainerName() {
				agent.setClient(ap.dialer(agent.config))
				agent.ready = true
				go agent.processTransactions()
				go agent.processBlocks()
			}
		}
	}
	ap.manageReadinessUnsafe()
	return nil
}

func (ap *AgentPool) handleStatusStopped(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	var newAgents []*Agent
	for _, agent := range ap.agents {
		var stopped bool
		for _, agentCfg := range payload {
			if agent.config.ContainerName() == agentCfg.ContainerName() {
				agent.Close()
				agent.ready = false
				stopped = true
				break
			}
		}
		if !stopped {
			newAgents = append(newAgents, agent)
		}
	}
	ap.agents = newAgents
	return nil
}

// manageReadinessUnsafe pauses or continues depending on the readiness.
func (ap *AgentPool) manageReadinessUnsafe() {
	var allReady bool
	for _, agent := range ap.agents {
		allReady = allReady || agent.ready
	}
	if allReady {
		processingState.Continue()
	} else {
		processingState.Pause()
	}
}

func (ap *AgentPool) registerMessageHandlers() {
	ap.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(ap.handleAgentVersionsUpdate))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusRunning, messaging.AgentsHandler(ap.handleStatusRunning))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusStopped, messaging.AgentsHandler(ap.handleStatusStopped))
}
