package containers

import (
	"fmt"

	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
)

// Label values
const (
	// LabelFortaSupervisor ensures that our docker client in the supervisor only touches
	// the containers managed by the supervisor service.
	LabelFortaSupervisor = "supervisor"
	LabelValueFortaIsBot = "true"
	// LabelValueStrategyVersion is for versioning the critical changes in container management strategy.
	// It's effective in deciding if a bot container should be re-created or not.
	LabelValueStrategyVersion = "2023-09-20T12:00:00Z"
)

// Limits define container limits.
type Limits struct {
	config.LogConfig
	config.BotResourceLimits
}

// NewBotContainerConfig creates a new bot container config.
func NewBotContainerConfig(
	networkID string, botConfig config.AgentConfig,
	logConfig config.LogConfig, resourcesConfig config.ResourcesConfig,
) docker.ContainerConfig {
	limits := config.GetAgentResourceLimits(resourcesConfig)

	return docker.ContainerConfig{
		Name:           botConfig.ContainerName(),
		Image:          botConfig.Image,
		NetworkID:      networkID,
		LinkNetworkIDs: []string{},
		Env: map[string]string{
			config.EnvJsonRpcHost:        config.DockerJSONRPCProxyContainerName,
			config.EnvJsonRpcPort:        config.DefaultJSONRPCProxyPort,
			config.EnvJWTProviderHost:    config.DockerJWTProviderContainerName,
			config.EnvJWTProviderPort:    config.DefaultJWTProviderPort,
			config.EnvPublicAPIProxyHost: config.DockerPublicAPIProxyContainerName,
			config.EnvPublicAPIProxyPort: config.DefaultPublicAPIProxyPort,
			config.EnvAgentGrpcPort:      botConfig.GrpcPort(),
			config.EnvFortaBotID:         botConfig.ID,
			config.EnvFortaBotOwner:      botConfig.Owner,
			config.EnvFortaChainID:       fmt.Sprintf("%d", botConfig.ChainID),
		},
		MaxLogFiles: logConfig.MaxLogFiles,
		MaxLogSize:  logConfig.MaxLogSize,
		CPUQuota:    limits.CPUQuota,
		Memory:      limits.Memory,
		Labels: map[string]string{
			docker.LabelFortaIsBot:                     LabelValueFortaIsBot,
			docker.LabelFortaSupervisorStrategyVersion: LabelValueStrategyVersion,
			docker.LabelFortaBotID:                     botConfig.ID,
		},
	}
}
