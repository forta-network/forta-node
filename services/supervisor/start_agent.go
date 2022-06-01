package supervisor

import (
	"errors"
	"fmt"
	"time"

	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"

	log "github.com/sirupsen/logrus"
)

var (
	errAgentAlreadyRunning = errors.New("agent already running")
)

func (sup *SupervisorService) startAgent(agent config.AgentConfig) error {
	if err := sup.agentImageClient.EnsureLocalImage(sup.ctx, fmt.Sprintf("agent %s", agent.ID), agent.Image); err != nil {
		return err
	}

	sup.mu.Lock()
	defer sup.mu.Unlock()

	_, ok := sup.getContainerUnsafe(agent.ContainerName())
	if ok {
		return errAgentAlreadyRunning
	}

	// run the bot admin first
	adminContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:  agent.AdminContainerName(),
		Image: sup.commonNodeImage,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "bot-admin"},
		Env: map[string]string{
			config.EnvContainerName: agent.AdminContainerName(), // consumed by server in the container
		},
		Volumes: map[string]string{
			sup.hostFortaDir: config.DefaultContainerFortaDirPath,
		},
		NetworkID:   sup.botNetworkID,
		MaxLogFiles: sup.maxLogFiles,
		MaxLogSize:  sup.maxLogSize,
		Labels: map[string]string{
			clients.DockerLabelFortaSupervisorStrategyVersion: SupervisorStrategyVersion,
		},
		Capabilities: []string{clients.DockerCapabilityNetAdmin}, // should be able to run iptables
	}, true)
	if err != nil {
		return err
	}

	// send rules request to the server
	log.WithField("container", agent.AdminContainerName()).Info("setting admin rules")
	for i := 0; i < 10; i++ {
		err = sup.botManager.SetBotAdminRules(agent.AdminContainerName())
		if err == nil {
			break
		}
		time.Sleep(time.Second * 3)
	}
	if err != nil {
		log.WithField("container", agent.AdminContainerName()).WithError(err).Info("failed to set admin rules")
		return err
	}
	log.WithField("container", agent.AdminContainerName()).Info("successfully set admin rules")

	limits := config.GetAgentResourceLimits(sup.config.Config.ResourcesConfig)

	agentContainer, err := sup.client.StartContainer(sup.ctx, clients.DockerContainerConfig{
		Name:      agent.ContainerName(),
		Image:     agent.Image,
		NetworkID: fmt.Sprintf("container:%s", agent.AdminContainerName()), // share network interfaces (and rules)
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

	sup.addContainerUnsafe(adminContainer, true)
	sup.addContainerUnsafe(agentContainer, false, &agent)

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

func (sup *SupervisorService) addContainerUnsafe(container *clients.DockerContainer, isAgentAdmin bool, agentConfig ...*config.AgentConfig) {
	if isAgentAdmin {
		sup.containers = append(sup.containers, &Container{
			DockerContainer: *container,
			IsAgentAdmin:    true,
		})
		return
	}
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
