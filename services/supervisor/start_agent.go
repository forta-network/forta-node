package supervisor

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/protocol"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/metrics"
	log "github.com/sirupsen/logrus"
)

var (
	errAgentAlreadyRunning = errors.New("agent already running")
)

const (
	agentStartTimeout = time.Minute * 5
)

func (sup *SupervisorService) emitMetric(agentID string, name string) {
	m := metrics.CreateAgentMetric(agentID, name, float64(1))
	sup.msgClient.PublishProto(
		messaging.SubjectMetricAgent,
		&protocol.AgentMetricList{Metrics: []*protocol.AgentMetric{m}},
	)
}

func (sup *SupervisorService) emitErrMetric(agentID string, name string, err error) {
	sup.emitMetric(agentID, fmt.Sprintf("%s-%s", name, err.Error()))
}

func (sup *SupervisorService) startAgent(ctx context.Context, agent config.AgentConfig) error {
	if err := sup.agentImageClient.EnsureLocalImage(ctx, fmt.Sprintf("agent %s", agent.ID), agent.Image); err != nil {
		sup.emitErrMetric(agent.ID, metrics.MetricAgentSupervisorStartErrorImage, err)
		return err
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()


	nwID, err := sup.client.CreatePublicNetwork(ctx, agent.ContainerName())
	if err != nil {
		sup.emitErrMetric(agent.ID, metrics.MetricAgentSupervisorStartErrorCreateNetwork, err)
		return err
	}

	_, ok := sup.getContainerUnsafe(agent.ContainerName())
	if ok {
		return errAgentAlreadyRunning
	}

	limits := config.GetAgentResourceLimits(sup.config.Config.ResourcesConfig)

	agentContainer, err := sup.client.StartContainer(
		ctx, clients.DockerContainerConfig{
			Name:           agent.ContainerName(),
			Image:          agent.Image,
			NetworkID:      nwID,
			LinkNetworkIDs: []string{},
			Env: map[string]string{
				config.EnvJsonRpcHost:        config.DockerJSONRPCProxyContainerName,
				config.EnvJsonRpcPort:        config.DefaultJSONRPCProxyPort,
				config.EnvJWTProviderHost:    config.DockerJWTProviderContainerName,
				config.EnvJWTProviderPort:    config.DefaultJWTProviderPort,
				config.EnvPublicAPIProxyHost: config.DockerPublicAPIProxyContainerName,
				config.EnvPublicAPIProxyPort: config.DefaultPublicAPIProxyPort,
				config.EnvAgentGrpcPort:      agent.GrpcPort(),
				config.EnvFortaBotID:         agent.ID,
				config.EnvFortaBotOwner:      agent.Owner,
				config.EnvFortaChainID:       fmt.Sprintf("%d", agent.ChainID),
			},
			MaxLogFiles: sup.maxLogFiles,
			MaxLogSize:  sup.maxLogSize,
			CPUQuota:    limits.CPUQuota,
			Memory:      limits.Memory,
			Labels: map[string]string{
				clients.DockerLabelFortaSupervisorStrategyVersion: SupervisorStrategyVersion,
			},
		},
	)
	if err != nil {
		sup.emitErrMetric(agent.ID, metrics.MetricAgentSupervisorErrorStartContainer, err)
		return err
	}

	// Attach the scanner, JWT Provider, Public API Proxy and the JSON-RPC proxy to the agent's network.
	for _, container := range sup.getBotNetworkContainers() {
		err := sup.client.AttachNetwork(ctx, container.ID, nwID)
		if err != nil {
			sup.emitErrMetric(agent.ID, metrics.MetricAgentSupervisorStartErrorAttachNetwork, err)
			return err
		}
	}

	sup.addContainerUnsafe(agentContainer, &agent)
	sup.emitMetric(agent.ID, metrics.MetricAgentSupervisorStartComplete)

	return nil
}

func (sup *SupervisorService) getContainerUnsafe(name string) (*Container, bool) {
	for _, container := range sup.containers {
		if container.Name == name {
			return container, true
		}
	}
	return nil, false
}

func (sup *SupervisorService) addContainerUnsafe(container *clients.DockerContainer, agentConfig ...*config.AgentConfig) {
	if agentConfig != nil {
		sup.containers = append(
			sup.containers, &Container{
				DockerContainer: *container,
				IsAgent:         true,
				AgentConfig:     agentConfig[0],
			},
		)
		return
	}
	sup.containers = append(sup.containers, &Container{DockerContainer: *container})
}

func (sup *SupervisorService) handleAgentRun(payload messaging.AgentPayload) error {
	return sup.handleAgentRunWithContext(sup.ctx, payload)
}

func (sup *SupervisorService) handleAgentRunWithContext(ctx context.Context, payload messaging.AgentPayload) error {
	sup.lastRun.Set()

	log.WithFields(
		log.Fields{
			"payload": len(payload),
		},
	).Infof("handle agent run")

	var wg sync.WaitGroup

	wg.Add(len(payload))

	for _, agent := range payload {
		go sup.doStartAgent(ctx, agent, &wg)
	}

	wg.Wait()

	return nil
}

// doStartAgent intended to use during multiple agent starts
func (sup *SupervisorService) doStartAgent(ctx context.Context, agent config.AgentConfig, wg *sync.WaitGroup) {
	sup.emitMetric(agent.ID, metrics.MetricAgentSupervisorStartBegin)

	ctx, cancel := context.WithTimeout(ctx, agentStartTimeout)
	defer cancel()

	defer wg.Done()

	logger := agentLogger(agent)

	err := sup.startAgent(ctx, agent)
	if err == errAgentAlreadyRunning {
		logger.Infof("agent container is already running - skipped")
		sup.emitMetric(agent.ID, metrics.MetricAgentSupervisorStartSkipAlreadyRunning)
		sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
		metrics.SendAgentMetrics(sup.msgClient, []*protocol.AgentMetric{metrics.CreateAgentMetric(agent.ID, metrics.MetricStatusRunning, 1)})
		return
	}
	if err != nil {
		logger.WithError(err).Error("failed to start agent")
		sup.emitErrMetric(agent.ID, metrics.MetricAgentSupervisorStartError, err)
		sup.msgClient.Publish(messaging.SubjectAgentsStatusStopped, messaging.AgentPayload{agent})
		return
	}

	// remove older versions of the agent
	containers, err := sup.client.GetContainers(ctx)
	for _, container := range containers {
		botContainerPrefix := fmt.Sprintf("%s-agent-%s", config.ContainerNamePrefix, utils.ShortenString(agent.ID, 8))
		containerName := container.Names[0]
		// check to see if there are bots with the same id, but different image
		if strings.HasPrefix(containerName, botContainerPrefix) && containerName != agent.ContainerName() {
			// emit stop action for outdated bot containers
			sup.msgClient.Publish(messaging.SubjectAgentsActionStop, messaging.AgentPayload{agent})
		}
	}

	// Broadcast the agent status.
	sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
	metrics.SendAgentMetrics(sup.msgClient, []*protocol.AgentMetric{metrics.CreateAgentMetric(agent.ID, metrics.MetricStatusRunning, 1)})
}

func (sup *SupervisorService) handleAgentStop(payload messaging.AgentPayload) error {
	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.lastStop.Set()

	stopped := make(map[string]bool)

	for _, agentCfg := range payload {
		sup.emitMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopBegin)
		logger := agentLogger(agentCfg)

		container, ok := sup.getContainerUnsafe(agentCfg.ContainerName())
		if !ok {
			sup.emitMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopSkipNotFound)
			logger.Warnf("container for agent was not found - skipping stop action")
			continue
		}
		if err := sup.client.StopContainer(sup.ctx, container.ID); err != nil {
			sup.emitErrMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopErrorContainer, err)
			return fmt.Errorf("failed to stop container '%s': %v", container.ID, err)
		}
		logger.Infof("successfully stopped the container")

		for _, c := range sup.getBotNetworkContainers() {
			if err := sup.client.DetachNetwork(sup.ctx, c.ID, agentCfg.ContainerName()); err != nil {
				sup.emitErrMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopErrorContainer, err)
				logger.WithError(err).Warnf("failed to disconnect container %s from network: %s", container.ID, agentCfg.ContainerName())
			}
		}

		if err := sup.client.RemoveNetworkByName(sup.ctx, agentCfg.ContainerName()); err != nil {
			sup.emitErrMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopErrorContainer, err)
			logger.WithError(err).Warnf("failed to remove container network: %s", agentCfg.ContainerName())
		}

		sup.emitMetric(agentCfg.ID, metrics.MetricAgentSupervisorStopComplete)
		stopped[container.ID] = true
	}

	// Remove the stopped agents from the list.
	var remainingContainers []*Container
	for _, container := range sup.containers {
		if !stopped[container.ID] {
			remainingContainers = append(remainingContainers, container)
		}
	}
	sup.containers = remainingContainers

	// Broadcast the agent statuses.
	if len(payload) > 0 {
		sup.msgClient.Publish(messaging.SubjectAgentsStatusStopped, payload)
	}
	return nil
}

func (sup *SupervisorService) registerMessageHandlers() {
	sup.msgClient.Subscribe(messaging.SubjectAgentsActionRun, messaging.AgentsHandler(sup.handleAgentRun))
	sup.msgClient.Subscribe(messaging.SubjectAgentsActionStop, messaging.AgentsHandler(sup.handleAgentStop))
	if sup.config.Config.InspectionConfig.InspectAtStartup {
		sup.msgClient.Subscribe(messaging.SubjectInspectionDone, messaging.InspectionResultsHandler(sup.handleInspectionResults))
	}
}

func (sup *SupervisorService) getBotNetworkContainers() []*clients.DockerContainer {
	return []*clients.DockerContainer{
		sup.scannerContainer, sup.jsonRpcContainer,
		sup.jwtProviderContainer, sup.publicAPIContainer,
	}
}

func agentLogger(agent config.AgentConfig) *log.Entry {
	return log.WithFields(
		log.Fields{
			"agentId": agent.ID, "image": agent.Image, "containerName": agent.ContainerName(),
		},
	)
}
