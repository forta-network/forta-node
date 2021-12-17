package updater

import (
	"context"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
	"os"
)

// Updater receives and starts the latest supervisor.
type Updater struct {
	ctx          context.Context
	cfg          config.Config
	imgStore     store.FortaImageStore
	dockerClient clients.DockerClient

	supervisorContainer *clients.DockerContainer
}

// NewUpdater creates a new updater.
func NewUpdater(ctx context.Context, cfg config.Config, imgStore store.FortaImageStore, dockerClient clients.DockerClient) *Updater {
	return &Updater{
		ctx:          ctx,
		cfg:          cfg,
		imgStore:     imgStore,
		dockerClient: dockerClient,
	}
}

// Start starts the service.
func (up *Updater) Start() error {
	go up.receive()
	return nil
}

// Name returns the name of the service.
func (up *Updater) Name() string {
	return "updater"
}

// Stop stops the service
func (up *Updater) Stop() error {
	up.stopSupervisor()
	return nil
}

func (up *Updater) stopSupervisor() {
	if up.supervisorContainer != nil {
		logger := log.WithField("container", up.supervisorContainer.ID)
		logger.Info("stopping supervisor")
		if err := up.dockerClient.InterruptContainer(context.Background(), up.supervisorContainer.ID); err != nil {
			logger.WithError(err).Error("error stopping supervisor container")
		} else {
			logger.Info("supervisor stopped")
		}
		if err := up.dockerClient.WaitContainerExit(context.Background(), up.supervisorContainer.ID); err != nil {
			logger.WithError(err).Error("error while waiting for container exist")
		}
	}
}

func (up *Updater) receive() {
	for latestRef := range up.imgStore.Latest() {
		log.WithField("image", latestRef).Info("received new node image reference")
		up.replaceSupervisor(latestRef)
	}
}

func (up *Updater) replaceSupervisor(imageRef string) {
	logger := log.WithField("image", imageRef)

	up.stopSupervisor()

	// ensure that we restart from scratch
	if err := up.dockerClient.Prune(up.ctx); err != nil {
		logger.WithError(err).Error("error while pruning after stopping old supervisor")
		return
	}
	if up.supervisorContainer != nil {
		if err := up.dockerClient.WaitContainerPrune(up.ctx, up.supervisorContainer.ID); err != nil {
			logger.WithError(err).Error("error while waiting for old supervisor prune")
			return
		}
	}

	if err := up.dockerClient.EnsureLocalImage(up.ctx, "supervisor", imageRef); err != nil {
		logger.WithError(err).Error("failed to ensure local image for supervisor")
		return
	}

	var err error

	log.Infof("forta-dir: %s", up.cfg.FortaDir)
	log.Infof("forta-config-path: %s", up.cfg.ConfigPath)

	cfgBytes, err := os.ReadFile(up.cfg.ConfigPath)
	if err != nil {
		return
	}
	up.supervisorContainer, err = up.dockerClient.StartContainer(up.ctx, clients.DockerContainerConfig{
		Name:  config.DockerSupervisorContainerName,
		Image: imageRef,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "supervisor"},
		Env: map[string]string{
			config.EnvConfigPath: up.cfg.ConfigPath,
			config.EnvNatsHost:   config.DockerNatsContainerName,
			config.EnvFortaDir:   up.cfg.FortaDir,
		},
		Volumes: map[string]string{
			"/var/run/docker.sock": "/var/run/docker.sock", // give access to host docker
			up.cfg.FortaDir:        config.DefaultContainerFortaDirPath,
		},
		Files: map[string][]byte{
			"passphrase":                      []byte(up.cfg.Passphrase),
			config.DefaultContainerConfigPath: cfgBytes,
		},
		MaxLogSize:  up.cfg.Log.MaxLogSize,
		MaxLogFiles: up.cfg.Log.MaxLogFiles,
	})
	if err != nil {
		logger.WithError(err).Errorf("failed to start the supervisor")
		return
	}
}
