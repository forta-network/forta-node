package agentpool

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
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
	ctx                     context.Context
	cfg                     config.Config
	agents                  []*poolagent.Agent
	txResults               chan *scanner.TxResult
	blockResults            chan *scanner.BlockResult
	combinationAlertResults chan *scanner.CombinationAlertResult
	msgClient               clients.MessageClient
	dialer                  func(config.AgentConfig) (clients.AgentClient, error)
	mu                      sync.RWMutex
	botWaitGroup            *sync.WaitGroup
	botReports              []*health.Report
}
// NewAgentPool creates a new agent pool.
func NewAgentPool(ctx context.Context, cfg config.Config, msgClient clients.MessageClient, waitBots int) *AgentPool {
	agentPool := &AgentPool{
		ctx:                     ctx,
		cfg:                     cfg,
		txResults:               make(chan *scanner.TxResult),
		blockResults:            make(chan *scanner.BlockResult),
		combinationAlertResults: make(chan *scanner.CombinationAlertResult),
		msgClient:               msgClient,
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

	reports := health.Reports{
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

	// add and reset existing reports
	for _, report := range ap.botReports {
		reports = append(reports, report)
	}

	ap.botReports = []*health.Report{}

	return reports
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

	if ap.botWaitGroup != nil {
		ap.botWaitGroup.Wait()
	}

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

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
		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending tx request to evalTxCh")

		// unblock req send and discard agent if agent is closed

		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
		case agent.TxRequestCh() <- &poolagent.TxRequest{
			Original: req,
			Encoded:  encoded,
		}:
		default: // do not try to send if the buffer is full
			lg.WithField("agent", agent.Config().ID).Debug("agent tx request buffer is full - skipping")
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

	if ap.botWaitGroup != nil {
		ap.botWaitGroup.Wait()
	}

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

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

		lg.WithFields(log.Fields{
			"agent":    agent.Config().ID,
			"duration": time.Since(startTime),
		}).Debug("sending block request to evalBlockCh")

		// unblock req send if agent is closed
		select {
		case <-agent.Closed():
			ap.discardAgent(agent)
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
	ap.msgClient.Publish(messaging.SubjectScannerBlock, &messaging.ScannerPayload{
		LatestBlockInput: blockNumber,
	})

	metrics.SendAgentMetrics(ap.msgClient, metricsList)
	lg.WithFields(log.Fields{
		"duration": time.Since(startTime),
	}).Debug("Finished SendEvaluateBlockRequest")
}

// SendEvaluateAlertRequest sends the request to all the active agents which
// should be processing the alert.
func (ap *AgentPool) SendEvaluateAlertRequest(req *protocol.EvaluateAlertRequest) {
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

	if ap.botWaitGroup != nil {
		ap.botWaitGroup.Wait()
	}

	ap.mu.RLock()
	agents := ap.agents
	ap.mu.RUnlock()

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
	case <-target.Closed():
		ap.discardAgent(target)
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

	ap.msgClient.Publish(messaging.SubjectScannerAlert, &messaging.ScannerPayload{})
	metrics.SendAgentMetrics(ap.msgClient, metricsList)
	lg.WithFields(
		log.Fields{
			"duration": time.Since(startTime),
		},
	).Debug("Finished SendEvaluateAlertRequest")
}

// CombinationAlertResults returns the receive-only alert results channel.
func (ap *AgentPool) CombinationAlertResults() <-chan *scanner.CombinationAlertResult {
	return ap.combinationAlertResults
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
	newAgents, agentsToRun := ap.updateAgentsOrFindMissing(latestVersions)

	// Find agents that are already deployed but doesn't exist in the latest versions payload
	agentsToStop := ap.findMissingAgentsInLatestVersions(latestVersions)

	ap.agents = newAgents

	ap.publishActions(agentsToRun, nil, agentsToStop, nil, nil)

	return nil
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

	ap.publishActions(nil, agentsReady, agentsToStop, newSubscriptions, removedSubscriptions)

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
	return nil
}

// updateAgentsOrFindMissing updates existing agents in the pool with the latest configuration
// and finds the missing agents in the pool to add them as new agents.
// Returns the updated list of agents and the configurations for the missing agents to start.
func (ap *AgentPool) updateAgentsOrFindMissing(latestVersions messaging.AgentPayload) ([]*poolagent.Agent, []config.AgentConfig) {
	var agents []*poolagent.Agent
	var agentsToRun []config.AgentConfig

	for _, agentCfg := range latestVersions {
		err := ap.findAgentAndHandle(
			agentCfg, func(agent *poolagent.Agent, _ *log.Entry) error {
				agent.SetShardConfig(agentCfg)
				agents = append(agents, agent)
				return nil
			},
		)
		if err != nil {
			// If the agent is missing in the pool, add it as a new agent.
			newAgent := poolagent.New(ap.ctx, agentCfg, ap.msgClient, ap.txResults, ap.blockResults, ap.combinationAlertResults)
			agents = append(agents, newAgent)
			agentsToRun = append(agentsToRun, agentCfg)
			log.WithField("agent", agentCfg.ID).Info("will trigger start")
		}
	}

	return agents, agentsToRun
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
	agentsToRun, agentsReady, agentsToStop []config.AgentConfig, newSubscriptions, removedSubscriptions []domain.CombinerBotSubscription,
) {
	if len(agentsToRun) > 0 {
		for _, agentConfig := range agentsToRun {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsActionRun,
					Status:  health.StatusInfo,
					Details: agentConfig.ID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsActionRun, agentsToRun)
	}
	if len(agentsReady) > 0 {
		for _, agentConfig := range agentsReady {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsStatusAttached,
					Status:  health.StatusInfo,
					Details: agentConfig.ID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsStatusAttached, agentsToRun)
	}
	if len(agentsToStop) > 0 {
		for _, agentConfig := range agentsToStop {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsActionStop,
					Status:  health.StatusInfo,
					Details: agentConfig.ID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsActionStop, agentsToStop)
	}
	if len(newSubscriptions) > 0 {
		for _, subscription := range newSubscriptions {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsAlertSubscribe,
					Status:  health.StatusInfo,
					Details: subscription.Subscriber.BotID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsAlertSubscribe, newSubscriptions)
	}
	if len(removedSubscriptions) > 0 {
		for _, subscription := range removedSubscriptions {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsAlertUnsubscribe,
					Status:  health.StatusInfo,
					Details: subscription.Subscriber.BotID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsAlertUnsubscribe, removedSubscriptions)
	}

	if len(agentsToRun) > 0 && ap.cfg.LocalModeConfig.IsStandalone() {
		for _, agentConfig := range agentsToRun {
			ap.botReports = append(
				ap.botReports, health.Report{
					Name:    messaging.SubjectAgentsStatusRunning,
					Status:  health.StatusInfo,
					Details: agentConfig.ID,
				},
			)
		}
		ap.msgClient.Publish(messaging.SubjectAgentsStatusRunning, agentsToRun)
	}
}

func (ap *AgentPool) registerMessageHandlers() {
	ap.msgClient.Subscribe(messaging.SubjectAgentsVersionsLatest, messaging.AgentsHandler(ap.handleAgentVersionsUpdate))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusRunning, messaging.AgentsHandler(ap.handleStatusRunning))
	ap.msgClient.Subscribe(messaging.SubjectAgentsStatusStopped, messaging.AgentsHandler(ap.handleStatusStopped))
}
