package agentpool

import (
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/protocol"
	"sync"
)

// TxResult contains request and response data.
type TxResult struct {
	Agent    Agent
	Request  *protocol.EvaluateTxRequest
	Response *protocol.EvaluateTxResponse
}

// BlockResult contains request and response data.
type BlockResult struct {
	Agent    Agent
	Request  *protocol.EvaluateBlockRequest
	Response *protocol.EvaluateBlockResponse
}

// AgentPool maintains the pool of agents that the scanner should
// interact with.
type AgentPool struct {
	agents       []Agent
	txResults    chan *TxResult
	blockResults chan *BlockResult
	mu           sync.RWMutex
}

// NewAgentPool creates a new agent pool.
func NewAgentPool() *AgentPool {
	return &AgentPool{
		txResults:    make(chan *TxResult, 100),
		blockResults: make(chan *BlockResult, 100),
	}
}

// SendEvaluateTxRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateTxRequest(req *protocol.EvaluateTxRequest) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	for _, agent := range ap.agents {
		if agent.shouldProcessBlock(req.Event.Block.BlockNumber) {
			agent.evalTxCh <- req
		}
	}
}

// TxResults returns the receive-only tx results channel.
func (ap *AgentPool) TxResults() <-chan *TxResult {
	return ap.txResults
}

// SendEvaluateBlockRequest sends the request to all of the active agents which
// should be processing the block.
func (ap *AgentPool) SendEvaluateBlockRequest(req *protocol.EvaluateBlockRequest) {
	ap.mu.RLock()
	defer ap.mu.RUnlock()
	for _, agent := range ap.agents {
		if agent.shouldProcessBlock(req.Event.BlockNumber) {
			agent.evalBlockCh <- req
		}
	}
}

// BlockResults returns the receive-only tx results channel.
func (ap *AgentPool) BlockResults() <-chan *BlockResult {
	return ap.blockResults
}

func (ap *AgentPool) handleAgentVersionsUpdate(payload interface{}) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	latestVersions := payload.([]config.AgentConfig)

	// The agents list which we completely replace with the old ones.
	var newAgents []Agent

	// Find the missing agents in the pool, add them to the new agents list
	// and send a "run" message.
	var agentsToRun []config.AgentConfig
	for _, agentCfg := range latestVersions {
		var found bool
		for _, agent := range ap.agents {
			found = found || (agent.config.Name == agentCfg.Name)
		}
		if !found {
			newAgents = append(newAgents, NewAgent(agentCfg, ap.txResults, ap.blockResults))
			agentsToRun = append(agentsToRun, agentCfg)
		}
	}
	if len(agentsToRun) > 0 {
		messaging.DefaultClient().Publish(messaging.SubjectAgentsActionRun, agentsToRun)
	}

	// Find the missing agents in the latest versions and send a "stop" message.
	// Otherwise, add to the new agents list so we keep on running.
	var agentsToStop []config.AgentConfig
	for _, agent := range ap.agents {
		var found bool
		var agentCfg config.AgentConfig
		for _, agentCfg = range latestVersions {
			found = found || (agent.config.Name == agentCfg.Name)
			if found {
				break
			}
		}
		if !found {
			agentsToStop = append(agentsToStop, agentCfg)
		} else {
			newAgents = append(newAgents, agent)
		}
	}
	if len(agentsToStop) > 0 {
		messaging.DefaultClient().Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}

	ap.agents = newAgents
	ap.manageReadinessUnsafe()
	return nil
}

func (ap *AgentPool) handleStatusRunning(payload interface{}) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	// If an agent was added before and just started to run, we should mark as ready
	// and start the processing goroutines.
	for _, agentCfg := range payload.([]config.AgentConfig) {
		for _, agent := range ap.agents {
			if agent.config.Name == agentCfg.Name {
				if err := agent.connect(); err != nil {
					return err
				}
				agent.ready = true
				go agent.processTransactions()
				go agent.processBlocks()
			}
		}
	}
	ap.manageReadinessUnsafe()
	return nil
}

func (ap *AgentPool) handleStatusStopped(payload interface{}) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()
	var newAgents []Agent
	for _, agent := range ap.agents {
		var stopped bool
		for _, agentCfg := range payload.([]config.AgentConfig) {
			if agent.config.Name == agentCfg.Name {
				stopped = true
				agent.Close()
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
	messaging.Subscribe(messaging.SubjectAgentsStatusRunning, ap.handleStatusRunning)
	messaging.Subscribe(messaging.SubjectAgentsStatusStopped, ap.handleStatusStopped)
}
