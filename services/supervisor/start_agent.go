package supervisor

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
)

var (
	errAgentAlreadyRunning = errors.New("agent already running")
)

const (
	agentStartTimeout = time.Minute * 5
)

func (sup *SupervisorService) startAgent(ctx context.Context, agent config.AgentConfig) error {
	if err := sup.agentImageClient.EnsureLocalImage(ctx, fmt.Sprintf("agent %s", agent.ID), agent.Image); err != nil {
		return err
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()

	_, ok := sup.getContainerUnsafe(agent.ContainerName())
	if ok {
		return errAgentAlreadyRunning
	}

	nwID, err := sup.client.CreatePublicNetwork(ctx, agent.ContainerName())
	if err != nil {
		return err
	}

	limits := config.GetAgentResourceLimits(sup.config.Config.ResourcesConfig)

	agentContainer, err := sup.client.StartContainer(
		ctx, clients.DockerContainerConfig{
			Name:           agent.ContainerName(),
			Image:          agent.Image,
			NetworkID:      nwID,
			LinkNetworkIDs: []string{},
			Env: map[string]string{
				config.EnvJsonRpcHost:     config.DockerJSONRPCProxyContainerName,
				config.EnvJsonRpcPort:     config.DefaultJSONRPCProxyPort,
				config.EnvJWTProviderHost: config.DockerJWTProviderContainerName,
				config.EnvJWTProviderPort: config.DefaultJWTProviderPort,
				config.EnvAgentGrpcPort:   agent.GrpcPort(),
				config.EnvFortaBotID:      agent.ID,
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
		return err
	}
	// Attach the scanner, JWT Provider and the JSON-RPC proxy to the agent's network.
	for _, containerID := range []string{
		sup.scannerContainer.ID, sup.jsonRpcContainer.ID,
		sup.jwtProviderContainer.ID,
	} {
		err := sup.client.AttachNetwork(ctx, containerID, nwID)
		if err != nil {
			return err
		}
	}

	sup.addContainerUnsafe(agentContainer, &agent)

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
	startCtx, cancel := context.WithTimeout(sup.ctx, agentStartTimeout)
	defer cancel()

	return sup.handleAgentRunWithContext(startCtx, payload)
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
	defer wg.Done()

	logger := agentLogger(agent)

	err := sup.startAgent(ctx, agent)
	if err == errAgentAlreadyRunning {
		logger.Infof("agent container is already running - skipped")
		sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
		return
	}
	if err != nil {
		logger.WithError(err).Error("failed to start agent")
		return
	}

	// Broadcast the agent status.
	sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
}

func (sup *SupervisorService) handleAgentStop(payload messaging.AgentPayload) error {
	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.lastStop.Set()

	stopped := make(map[string]bool)
	for _, agentCfg := range payload {
		logger := agentLogger(agentCfg)

		container, ok := sup.getContainerUnsafe(agentCfg.ContainerName())
		if !ok {
			logger.Warnf("container for agent was not found - skipping stop action")
			continue
		}
		if err := sup.client.StopContainer(sup.ctx, container.ID); err != nil {
			return fmt.Errorf("failed to stop container '%s': %v", container.ID, err)
		}
		logger.Infof("successfully stopped the container")
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
}

func agentLogger(agent config.AgentConfig) *log.Entry {
	return log.WithFields(
		log.Fields{
			"agentId": agent.ID, "image": agent.Image, "containerName": agent.ContainerName(),
		},
	)
}
