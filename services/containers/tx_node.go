package containers

import (
	"context"
	"encoding/json"
	"fmt"
	"path"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
)

// TxNodeService manages the safe-node docker container as a service
type TxNodeService struct {
	ctx        context.Context
	client     clients.DockerClient
	authClient clients.DockerClient

	msgClient   clients.MessageClient
	config      TxNodeServiceConfig
	maxLogSize  string
	maxLogFiles int

	scannerContainer *clients.DockerContainer
	jsonRpcContainer *clients.DockerContainer
	containers       []*clients.DockerContainer
	mu               sync.RWMutex
}

type TxNodeServiceConfig struct {
	Config     config.Config
	Passphrase string
}

func (t *TxNodeService) Start() error {
	if err := t.start(); err != nil {
		return err
	}

	t.msgClient = messaging.NewClient("cli", ":"+config.DefaultNatsPort) // accessible from localhost
	t.registerMessageHandlers()

	go t.healthCheck()

	return nil
}

func (t *TxNodeService) start() error {
	log.Infof("Starting %s", t.Name())
	_, err := log.ParseLevel(t.config.Config.Log.Level)
	if err != nil {
		log.Error("invalid log level", err)
		return err
	}

	t.maxLogSize = t.config.Config.Log.MaxLogSize
	t.maxLogFiles = t.config.Config.Log.MaxLogFiles

	cfgBytes, err := json.Marshal(t.config.Config)
	if err != nil {
		log.Error("cannot marshal config to json", err)
		return err
	}
	cfgJson := string(cfgBytes)

	alertsDBPath := t.config.Config.Query.DB.Path
	if len(alertsDBPath) == 0 {
		alertsDBPath = path.Join(t.config.Config.FortaDir, "alertsdb")
	}

	if err := t.client.Prune(t.ctx); err != nil {
		return err
	}

	if config.UseDockerImages == "remote" {
		if err := t.ensureNodeImages(); err != nil {
			return err
		}
	}

	nodeNetwork, err := t.client.CreatePublicNetwork(t.ctx, config.DockerNetworkName)
	if err != nil {
		return err
	}

	natsContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerNatsContainerName,
		Image: "nats:2.3.2",
		Ports: map[string]string{
			"4222": "4222",
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	queryContainerVolumes := map[string]string{
		alertsDBPath:             store.DBPath,
		t.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
	}
	// Mount the test alerts dir only if we really should write test alerts to a file.
	testAlertsCfg := t.config.Config.Query.PublishTo.TestAlerts
	if !testAlertsCfg.Disable && len(testAlertsCfg.WebhookURL) == 0 && len(t.config.Config.LocalAgents) > 0 {
		queryContainerVolumes["/tmp"] = "/test-alerts"
	}
	queryContainer, err := t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerQueryContainerName,
		Image: config.DockerQueryContainerImage,
		Env: map[string]string{
			config.EnvConfig:   cfgJson,
			config.EnvFortaDir: config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			fmt.Sprintf("%d", t.config.Config.Query.Port): "80",
		},
		Volumes: queryContainerVolumes,
		Files: map[string][]byte{
			"passphrase": []byte(t.config.Passphrase),
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.jsonRpcContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerJSONRPCProxyContainerName,
		Image: config.DockerJSONRPCProxyContainerImage,
		Env: map[string]string{
			config.EnvConfig: cfgJson,
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.scannerContainer, err = t.client.StartContainer(t.ctx, clients.DockerContainerConfig{
		Name:  config.DockerScannerContainerName,
		Image: config.DockerScannerContainerImage,
		Env: map[string]string{
			config.EnvConfig:    cfgJson,
			config.EnvFortaDir:  config.DefaultContainerFortaDirPath,
			config.EnvQueryNode: config.DockerQueryContainerName,
			config.EnvNatsHost:  config.DockerNatsContainerName,
		},
		Volumes: map[string]string{
			t.config.Config.FortaDir: config.DefaultContainerFortaDirPath,
		},
		Files: map[string][]byte{
			"passphrase": []byte(t.config.Passphrase),
		},
		NetworkID:   nodeNetwork,
		MaxLogFiles: t.maxLogFiles,
		MaxLogSize:  t.maxLogSize,
	})
	if err != nil {
		return err
	}

	t.addContainerUnsafe(natsContainer, queryContainer, t.jsonRpcContainer, t.scannerContainer)

	return nil
}

func (t *TxNodeService) ensureNodeImages() error {
	for _, image := range []struct {
		Name        string
		Ref         string
		RequireAuth bool
	}{
		{
			Name: "nats",
			Ref:  "nats:2.3.2",
		},
		{
			Name:        "scanner",
			Ref:         config.DockerScannerContainerImage,
			RequireAuth: true,
		},
		{
			Name:        "query",
			Ref:         config.DockerQueryContainerImage,
			RequireAuth: true,
		},
		{
			Name:        "json-rpc",
			Ref:         config.DockerJSONRPCProxyContainerImage,
			RequireAuth: true,
		},
	} {
		if err := t.ensureLocalImage(image.Name, image.Ref, image.RequireAuth); err != nil {
			return err
		}
	}
	return nil
}

func (t *TxNodeService) ensureLocalImage(name, ref string, requireAuth bool) error {
	if t.client.HasLocalImage(t.ctx, ref) {
		log.Infof("found local image for '%s': %s", name, ref)
		return nil
	}
	client := t.client
	if requireAuth {
		client = t.authClient
	}
	err := client.PullImage(t.ctx, ref)
	if err != nil {
		return fmt.Errorf("failed to pull image (%s, %s): %v", name, ref, err)
	}
	log.Infof("pulled image for '%s': %s", name, ref)
	return nil
}

func (t *TxNodeService) Stop() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	ctx := context.Background()
	for _, cnt := range t.containers {
		if err := t.client.StopContainer(ctx, cnt.ID); err != nil {
			log.Error(fmt.Sprintf("error stopping %s container", cnt.ID), err)
		} else {
			log.Infof("Container %s is stopped", cnt.ID)
		}
	}
	return nil
}

func (t *TxNodeService) Name() string {
	return "TxNode"
}

func NewTxNodeService(ctx context.Context, cfg TxNodeServiceConfig) (*TxNodeService, error) {
	dockerAuthClient, err := clients.NewAuthDockerClient(cfg.Config.Registry.Username, cfg.Config.Registry.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to create the agent docker client: %v", err)
	}
	dockerClient, err := clients.NewDockerClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create the docker client: %v", err)
	}
	return &TxNodeService{
		ctx:        ctx,
		client:     dockerClient,
		authClient: dockerAuthClient,
		config:     cfg,
	}, nil
}
