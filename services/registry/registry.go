package registry

import (
	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/services"
)

// RegistryService listens to the agent pool changes so the node
// can stay in sync.
// TODO: Instead of publishing messages for the static config,
// check or listen to an actively maintained resource (e.g. a smart contract).
// TODO: The registry service or the config should construct unique names for the agents.
type RegistryService struct {
	cfg       config.Config
	msgClient clients.MessageClient
}

// New creates a new service.
func New(cfg config.Config, msgClient clients.MessageClient) services.Service {
	return &RegistryService{
		cfg:       cfg,
		msgClient: msgClient,
	}
}

// Start starts the registry service.
func (rs *RegistryService) Start() error {
	rs.msgClient.Publish(messaging.SubjectAgentsVersionsLatest, rs.cfg.Agents)
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
