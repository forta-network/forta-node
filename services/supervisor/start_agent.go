package supervisor

import (
	"errors"
	"fmt"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
)

var (
	errAgentAlreadyRunning = errors.New("agent already running")
)

func (sup *SupervisorService) startAgent(agent config.AgentConfig) error {
	if err := sup.client.EnsureLocalImage(sup.ctx, fmt.Sprintf("agent %s", agent.ID), agent.Image); err != nil {
		return err
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()

	_, ok := sup.getContainerUnsafe(agent.ContainerName())
	if ok {
		return errAgentAlreadyRunning
	}

	nwID, err := sup.client.CreatePublicNetwork(sup.ctx, agent.ContainerName())
	if err != nil {
		return err
	}

	limits := config.GetAgentResourceLimits(sup.config.Config.ResourcesConfig)

	agentContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:           agent.ContainerName(),
		Image:          agent.Image,
		NetworkID:      nwID,
		LinkNetworkIDs: []string{},
		Env: map[string]string{
			config.EnvJsonRpcHost:   config.DockerJSONRPCProxyContainerName,
			config.EnvJsonRpcPort:   "8545",
			config.EnvAgentGrpcPort: agent.GrpcPort(),
		},
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
		CPUQuota:    limits.CPUQuota,
		Memory:      limits.Memory,
		Labels: map[string]string{
			clients.DockerLabelFortaSupervisorStrategyVersion: SupervisorStrategyVersion,
		},
	})
	if err != nil {
		return err
	}
	// Attach the scanner and the JSON-RPC proxy to the agent's network.
	for _, containerID := range []string{sup.scannerContainer.ID, sup.jsonRpcContainer.ID} {
		err := sup.client.AttachNetwork(sup.ctx, containerID, nwID)
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
		sup.containers = append(sup.containers, &Container{
			DockerContainer: *container,
			IsAgent:         true,
			AgentConfig:     agentConfig[0],
		})
		return
	}
	sup.containers = append(sup.containers, &Container{DockerContainer: *container})
}

func (sup *SupervisorService) handleAgentRun(payload messaging.AgentPayload) error {
	sup.lastRun.Set()

	log.WithFields(log.Fields{
		"payload": len(payload),
	}).Infof("handle agent run")

	for _, agent := range payload {
		err := sup.startAgent(agent)
		if err == errAgentAlreadyRunning {
			log.Infof("agent container '%s' is already running - skipped", agent.ContainerName())
			sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
			continue
		}
		if err != nil {
			log.Errorf("failed to start agent: %v", err)
			continue
		}

		// Broadcast the agent status.
		sup.msgClient.Publish(messaging.SubjectAgentsStatusRunning, messaging.AgentPayload{agent})
	}
	return nil
}

func (sup *SupervisorService) handleAgentStop(payload messaging.AgentPayload) error {
	sup.mu.Lock()
	defer sup.mu.Unlock()

	sup.lastStop.Set()

	stopped := make(map[string]bool)
	for _, agentCfg := range payload {
		container, ok := sup.getContainerUnsafe(agentCfg.ContainerName())
		if !ok {
			log.Warnf("container for agent '%s' was not found - skipping stop action", agentCfg.ContainerName())
			continue
		}
		if err := sup.client.StopContainer(sup.ctx, container.ID); err != nil {
			return fmt.Errorf("failed to stop container '%s': %v", container.ID, err)
		}
		log.Infof("successfully stopped the container: %v", agentCfg.ContainerName())
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
