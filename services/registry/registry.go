package registry

import (
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/services"
)

// RegistryService listens to the agent pool changes so the node
// can stay in sync.
// TODO: Instead of publishing messages for the static config,
// check or listen to an actively maintained resource (e.g. a smart contract).
type RegistryService struct {
	cfg config.Config
}

// New creates a new service.
func New(cfg config.Config) services.Service {
	return &RegistryService{
		cfg: cfg,
	}
}

// Start starts the registry service.
func (rs *RegistryService) Start() error {
	messaging.Publish(messaging.SubjectAgentsVersionsLatest, rs.cfg.Agents)
	return nil
}

// Stop stops the registry service.
func (rs *RegistryService) Stop() error {
	return nil
}

// Name returns the name of the service.
func (rs *RegistryService) Name() string {
	return "RegistryService"
}
