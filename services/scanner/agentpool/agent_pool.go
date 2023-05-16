package agentpool

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-node/metrics"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/agentgrpc"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/services/scanner"
	"github.com/forta-network/forta-node/services/scanner/agentpool/poolagent"
	log "github.com/sirupsen/logrus"
)

// AgentPool maintains the pool of agents that the scanner should
// interact with.
type AgentPool struct {
	ctx    context.Context
	cfg    config.Config
	agents []*poolagent.Agent

	msgClient    clients.MessageClient
	dialer       func(config.AgentConfig) (clients.AgentClient, error)
	mu           sync.RWMutex
	botWaitGroup *sync.WaitGroup
	dispatcher   *Dispatcher
}

// NewAgentPool creates a new agent pool.
func NewAgentPool(ctx context.Context, cfg config.Config, msgClient clients.MessageClient, waitBots int) *AgentPool {
	agentPool := &AgentPool{
		ctx: ctx,
		cfg: cfg,
		dispatcher: &Dispatcher{
			txResults:               make(chan *scanner.TxResult),
			blockResults:            make(chan *scanner.BlockResult),
			combinationAlertResults: make(chan *scanner.CombinationAlertResult),
			msgClient:               msgClient,
		},
		msgClient: msgClient,
		dialer: func(ac config.AgentConfig) (clients.AgentClient, error) {
			client := agentgrpc.NewClient()
			if err := client.Dial(ac); err != nil {
				return nil, err
			}
			return client, nil
		},
	}
	if waitBots > 0 {
		agentPool.botWaitGroup = &sync.WaitGroup{}
		agentPool.botWaitGroup.Add(waitBots)
		go agentPool.logBotWait()
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

func (ap *AgentPool) logBotWait() {
	if ap.botWaitGroup != nil {
		ap.botWaitGroup.Wait()
		log.Info("started all bots")
	}
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
	ap.dispatcher.SetAgents(newAgents)
	ap.mu.Unlock()
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

// handleAgentVersionsUpdate updates the list of agents in the AgentPool based on the latest
// versions provided in the payload. It finds the missing agents in the pool and adds them to a new
// agents list. Then, it sends a "run" message for the newly added agents. It also finds the missing
// agents in the latest versions and sends a "stop" message for those. The agents that already exist
// in the pool and the latest versions will remain in the pool. Finally, it publishes the "run" and "stop"
// actions along with any necessary subscription updates.
func (ap *AgentPool) handleAgentVersionsUpdate(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	latestVersions := payload

	// Find the missing agents in the pool, add them to the new agents list
	// and send a "run" message.
	// newAgents is the updated list of all agents
	// agentsToRun is the list of missing agents
	newAgents, agentsToRun, updatedAgents := ap.updateAgentsOrFindMissing(latestVersions)

	// Find agents that are already deployed but doesn't exist in the latest versions payload
	agentsToStop := ap.findMissingAgentsInLatestVersions(latestVersions)

	ap.agents = newAgents
	ap.dispatcher.SetAgents(newAgents)

	ap.publishActions(agentsToRun, nil, agentsToStop, updatedAgents, nil, nil)

	return nil
}
func (ap *AgentPool) SetAgents(payload []*poolagent.Agent) error {

}
func (ap *AgentPool) handleStatusRunning(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	// If an agent was added before and just started to run, we should mark as ready.
	var agentsToStop []config.AgentConfig
	var agentsReady []config.AgentConfig
	var newSubscriptions []domain.CombinerBotSubscription
	var removedSubscriptions []domain.CombinerBotSubscription

	for _, agentCfg := range payload {
		_ = ap.findAgentAndHandle(
			agentCfg, func(agent *poolagent.Agent, logger *log.Entry) error {
				if agent.IsReady() {
					return nil
				}

				c, err := ap.dialer(agent.Config())
				if err != nil {
					log.WithField("agent", agent.Config().ID).WithError(err).Error("handleStatusRunning: error while dialing")
					agentsToStop = append(agentsToStop, agent.Config())
					removedSubscriptions = append(removedSubscriptions, agent.CombinerBotSubscriptions()...)
					return nil
				}

				agent.SetClient(c)
				agent.SetReady()
				agent.StartProcessing()
				agent.WaitInitialization()

				newSubscriptions = append(newSubscriptions, agent.CombinerBotSubscriptions()...)

				logger.WithField("image", agent.Config().Image).Info("attached")
				agentsReady = append(agentsReady, agent.Config())
				return nil
			},
		)
	}

	ap.publishActions(nil, agentsReady, agentsToStop, nil, newSubscriptions, removedSubscriptions)

	if ap.botWaitGroup != nil && len(agentsReady) > 0 {
		ap.botWaitGroup.Add(-len(agentsReady))
	}

	return nil
}

func (ap *AgentPool) handleStatusStopped(payload messaging.AgentPayload) error {
	ap.mu.Lock()
	defer ap.mu.Unlock()

	var newAgents []*poolagent.Agent
	var removedSubscriptions []domain.CombinerBotSubscription
	for _, agent := range ap.agents {
		var stopped bool
		for _, agentCfg := range payload {
			if agent.Config().ContainerName() == agentCfg.ContainerName() {
				agent.Close()
				log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Info("detached")
				stopped = true
				removedSubscriptions = append(removedSubscriptions, agent.CombinerBotSubscriptions()...)
				break
			}
		}
		if !stopped {
			log.WithField("agent", agent.Config().ID).WithField("image", agent.Config().Image).Debug("not stopped")
			newAgents = append(newAgents, agent)
		}
	}

	if len(removedSubscriptions) > 0 {
		ap.msgClient.Publish(messaging.SubjectAgentsAlertUnsubscribe, removedSubscriptions)
	}

	ap.agents = newAgents
	ap.dispatcher.SetAgents(newAgents)
	return nil
}

// updateAgentsOrFindMissing updates existing agents in the pool with the latest configuration
// and finds the missing agents in the pool to add them as new agents.
// Returns the updated list of agents and the configurations for the missing agents to start.
func (ap *AgentPool) updateAgentsOrFindMissing(latestVersions messaging.AgentPayload) ([]*poolagent.Agent, []config.AgentConfig, []config.AgentConfig) {
	var agents []*poolagent.Agent
	var updatedAgents []config.AgentConfig
	var agentsToRun []config.AgentConfig

	for _, agentCfg := range latestVersions {
		err := ap.findAgentAndHandle(
			agentCfg, func(agent *poolagent.Agent, _ *log.Entry) error {
				if !agent.Config().Equal(agentCfg) {
					updatedAgents = append(updatedAgents, agent.Config())
				}
				agent.UpdateConfig(agentCfg)
				agents = append(agents, agent)
				return nil
			},
		)
		if err != nil {
			// If the agent is missing in the pool, add it as a new agent.
			newAgent := poolagent.New(ap.ctx, agentCfg, ap.msgClient, ap.dispatcher.txResults, ap.dispatcher.blockResults, ap.dispatcher.combinationAlertResults)
			agents = append(agents, newAgent)
			agentsToRun = append(agentsToRun, agentCfg)
			log.WithField("agent", agentCfg.ID).Info("will trigger start")
		}
	}

	return agents, agentsToRun, updatedAgents
}

// findMissingAgentsInLatestVersions finds agents in the pool that are not in the latest versions payload,
// and returns the list of these agents to stop.
func (ap *AgentPool) findMissingAgentsInLatestVersions(latestVersions messaging.AgentPayload) []config.AgentConfig {
	var agentsToStop []config.AgentConfig

	for _, agent := range ap.agents {
		cfg := agent.Config()
		found := false
		for _, agentCfg := range latestVersions {
			if agentCfg.ContainerName() == agent.Config().ContainerName() {
				found = true
				break
			}
		}

		if !found {
			_ = agent.Close()
			agentsToStop = append(agentsToStop, cfg)
			log.WithField("agent", cfg.ID).WithField("image", cfg.Image).Info("will trigger stop")
		}
	}

	return agentsToStop
}

func (ap *AgentPool) findAgentAndHandle(cfg config.AgentConfig, handler func(agent *poolagent.Agent, logger *log.Entry) error) error {
	for _, agent := range ap.agents {
		if cfg.ContainerName() == agent.Config().ContainerName() {
			logger := log.WithField("agent", agent.Config().ID)
			return handler(agent, logger)
		}
	}
	return fmt.Errorf("agent not found: %v", cfg.ID)
}

func (ap *AgentPool) publishActions(
	agentsToRun, agentsReady, agentsToStop, agentsUpdated []config.AgentConfig, newSubscriptions, removedSubscriptions []domain.CombinerBotSubscription,
) {
	var ms []*protocol.AgentMetric

	if len(agentsToRun) > 0 {
		for _, agentConfig := range agentsToRun {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricActionRun, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsActionRun, agentsToRun)
	}

	if len(agentsReady) > 0 {
		for _, agentConfig := range agentsReady {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricStatusAttached, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsStatusAttached, agentsToRun)
	}

	if len(agentsToStop) > 0 {
		for _, agentConfig := range agentsToStop {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricActionStop, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}

	if len(agentsUpdated) > 0 {
		for _, agentConfig := range agentsUpdated {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricActionUpdated, 1))
		}
	}

	if len(newSubscriptions) > 0 {
		for _, agentConfig := range newSubscriptions {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.Subscriber.BotID, metrics.MetricActionSubscribe, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsAlertSubscribe, newSubscriptions)
	}

	if len(removedSubscriptions) > 0 {
		for _, agentConfig := range removedSubscriptions {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.Subscriber.BotID, metrics.MetricActionUnsubscribe, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsAlertUnsubscribe, removedSubscriptions)
	}

	if len(agentsToRun) > 0 && ap.cfg.LocalModeConfig.IsStandalone() {
		for _, agentConfig := range agentsToRun {
			ms = append(ms, metrics.CreateAgentMetric(agentConfig.ID, metrics.MetricStatusRunning, 1))
		}
		ap.msgClient.Publish(messaging.SubjectAgentsStatusRunning, agentsToRun)
	}

	metrics.SendAgentMetrics(ap.msgClient, ms)
}

func (ap *AgentPool) registerMessageHandlers() {
	ap.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(ap.handleAgentVersionsUpdate))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusRunning, messaging.AgentsHandler(ap.handleStatusRunning))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusStopped, messaging.AgentsHandler(ap.handleStatusStopped))
}
