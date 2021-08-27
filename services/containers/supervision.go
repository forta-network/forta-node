package containers

import (
	"fmt"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
)

func (t *TxNodeService) startAgent(agent config.AgentConfig) error {
	if err := t.ensureLocalImage(fmt.Sprintf("agent %s", agent.ID), agent.Image, true); err != nil {
		return err
	}
	log.Infof("pulled agent image: %v", agent.Image)

	nwID, err := t.client.CreatePublicNetwork(t.ctx, agent.ContainerName())
	if err != nil {
		return err
	}
	agentContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:           agent.ContainerName(),
		Image:          agent.Image,
		NetworkID:      nwID,
		LinkNetworkIDs: []string{},
		Env: map[string]string{
			config.EnvJsonRpcHost:   config.DockerJSONRPCProxyContainerName,
			config.EnvJsonRpcPort:   "8545",
			config.EnvAgentGrpcPort: agent.GrpcPort(),
		},
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}
	// Attach the scanner and the JSON-RPC proxy to the agent's network.
	for _, containerID := range []string{t.scannerContainer.ID, t.jsonRpcContainer.ID} {
		err := t.client.AttachNetwork(t.ctx, containerID, nwID)
		if err == clients.ErrAlreadyExistsInNetwork {
			continue
		}
		if err != nil {
			return err
		}
	}

	t.addContainerUnsafe(agentContainer)

	return nil
}

func (t *TxNodeService) getContainerUnsafe(name string) (*clients.DockerContainer, bool) {
	for _, container := range t.containers {
		if container.Name == name {
			return container, true
		}
	}
	return nil, false
}

func (t *TxNodeService) addContainerUnsafe(container ...*clients.DockerContainer) {
	t.containers = append(t.containers, container...)
}

func (t *TxNodeService) handleAgentRun(payload messaging.AgentPayload) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	for _, agent := range payload {
		_, ok := t.getContainerUnsafe(agent.ContainerName())
		if ok {
			log.Infof("agent container '%s' is already running - skipped", agent.ContainerName())
			continue
		}

		if err := t.startAgent(agent); err != nil {
			return err
		}
	}
	// Broadcast the agent statuses.
	t.msgClient.Publish(messaging.SubjectAgentsStatusRunning, payload)
	return nil
}

func (t *TxNodeService) handleAgentStop(payload messaging.AgentPayload) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	stopped := make(map[string]bool)
	for _, agentCfg := range payload {
		container, ok := t.getContainerUnsafe(agentCfg.ContainerName())
		if !ok {
			log.Warnf("container for agent '%s' was not found - skipping stop action", agentCfg.ContainerName())
			continue
		}
		if err := t.client.StopContainer(t.ctx, container.ID); err != nil {
			return fmt.Errorf("failed to stop container '%s': %v", container.ID, err)
		}
		log.Infof("successfully stopped the container: %v", agentCfg.ContainerName())
		stopped[container.ID] = true
	}

	// Remove the stopped agents from the list.
	var remainingContainers []*clients.DockerContainer
	for _, container := range t.containers {
		if !stopped[container.ID] {
			remainingContainers = append(remainingContainers, container)
		}
	}
	t.containers = remainingContainers

	// Broadcast the agent statuses.
	if len(payload) > 0 {
		t.msgClient.Publish(messaging.SubjectAgentsStatusStopped, payload)
	}
	return nil
}

func (t *TxNodeService) registerMessageHandlers() {
	t.msgClient.Subscribe(messaging.SubjectAgentsActionRun, messaging.AgentsHandler(t.handleAgentRun))
	t.msgClient.Subscribe(messaging.SubjectAgentsActionStop, messaging.AgentsHandler(t.handleAgentStop))
}
