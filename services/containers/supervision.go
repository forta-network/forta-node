package containers

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func (t *TxNodeService) startAgent(agent config.AgentConfig) error {
	nwID, err := t.client.CreatePublicNetwork(t.ctx, agent.Name)
	if err != nil {
		return err
	}
	agentContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:           agent.ContainerName(),
		Image:          agent.Image,
		NetworkID:      nwID,
		LinkNetworkIDs: []string{},
		Env: map[string]string{
			config.EnvJsonRpcHost:   jsonRpcProxyName,
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
		if err := t.client.AttachNetwork(t.ctx, containerID, nwID); err != nil {
			return err
		}
	}

	t.addContainer(agentContainer)

	return nil
}

func (t *TxNodeService) getContainer(name string) (*clients.DockerContainer, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	for _, container := range t.containers {
		if container.Name == name {
			return container, true
		}
	}
	return nil, false
}

func (t *TxNodeService) addContainer(container ...*clients.DockerContainer) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.containers = append(t.containers, container...)
}

func (t *TxNodeService) handleAgentRun(payload interface{}) error {
	// TODO: Be careful about the agent names (container names) when running two
	// versions of the same agent. The registry should enforce the correct names.
	// TODO: Should be idempotent to handle message redelivery cases.
	for _, agent := range payload.(messaging.AgentPayload) {
		if err := t.startAgent(agent); err != nil {
			return err
		}
	}
	// Broadcast the agent statuses.
	messaging.Publish(messaging.SubjectAgentsStatusRunning, payload)
	return nil
}

func (t *TxNodeService) handleAgentStop(payload interface{}) error {
	t.mu.Lock()
	defer t.mu.RUnlock()

	stopped := make(map[string]bool)
	for _, agentCfg := range payload.(messaging.AgentPayload) {
		container, ok := t.getContainer(agentCfg.Name)
		if !ok {
			log.Warnf("container for agent '%s' was not found - skipping stop action", agentCfg.Name)
			continue
		}
		if err := t.client.StopContainer(t.ctx, container.ID); err != nil {
			return fmt.Errorf("failed to stop container '%s': %v", container.ID, err)
		}
		log.Infof("successfully stopped the container: %v", agentCfg.Name)
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
	messaging.Publish(messaging.SubjectAgentsStatusStopped, payload)
	return nil
}

func (tx *TxNodeService) registerMessageHandlers() {
	messaging.Subscribe(messaging.SubjectAgentsActionRun, tx.handleAgentRun)
	messaging.Subscribe(messaging.SubjectAgentsActionStop, tx.handleAgentStop)
}
