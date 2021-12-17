package runner

import (
	"context"
	"encoding/json"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	log "github.com/sirupsen/logrus"
)

// Runner receives and starts the latest updater and supervisor.
type Runner struct {
	ctx          context.Context
	cfg          config.Config
	imgStore     store.FortaImageStore
	dockerClient clients.DockerClient

	updaterPort         string
	updaterContainer    *clients.DockerContainer
	supervisorContainer *clients.DockerContainer
}

// NewRunner creates a new runner.
func NewRunner(ctx context.Context, cfg config.Config,
	imgStore store.FortaImageStore, dockerClient clients.DockerClient,
	updaterPort string,
) *Runner {
	return &Runner{
		ctx:          ctx,
		cfg:          cfg,
		imgStore:     imgStore,
		dockerClient: dockerClient,
		updaterPort:  updaterPort,
	}
}

// Start starts the service.
func (runner *Runner) Start() error {
	go runner.receive()
	return nil
}

// Name returns the name of the service.
func (runner *Runner) Name() string {
	return "runner"
}

// Stop stops the service
func (runner *Runner) Stop() error {
	runner.stopContainer(runner.updaterContainer)
	runner.stopContainer(runner.supervisorContainer)
	return nil
}

func (runner *Runner) stopContainer(container *clients.DockerContainer) {
	if container != nil {
		logger := log.WithField("container", container.ID).WithField("name", container.Name)
		logger.Info("stopping")
		if err := runner.dockerClient.InterruptContainer(context.Background(), container.ID); err != nil {
			logger.WithError(err).Error("error stopping container")
		} else {
			logger.Info("stopped")
		}
		if err := runner.dockerClient.WaitContainerExit(context.Background(), container.ID); err != nil {
			logger.WithError(err).Error("error while waiting for container exit")
		}
	}
}

func (runner *Runner) receive() {
	for latestRefs := range runner.imgStore.Latest() {
		log.WithField("supervisor", latestRefs.Supervisor).WithField("updater", latestRefs.Updater).
			Info("received new node image reference")
		runner.replaceContainers(latestRefs)
	}
}

func (runner *Runner) replaceContainers(imageRefs store.ImageRefs) {
	logger := log.WithField("supervisor", imageRefs.Supervisor).WithField("updater", imageRefs.Updater)
	if imageRefs.Release != nil {
		logger = logger.WithField("version", imageRefs.Release.Metadata.Version)
	}

	runner.Stop()

	// ensure that we restart from scratch
	if err := runner.dockerClient.Prune(runner.ctx); err != nil {
		logger.WithError(err).Error("error while pruning after stopping old containers")
		return
	}
	for _, container := range []*clients.DockerContainer{runner.updaterContainer, runner.supervisorContainer} {
		if container != nil {
			if err := runner.dockerClient.WaitContainerPrune(runner.ctx, container.ID); err != nil {
				logger.WithError(err).Error("error while waiting for old container prune")
				return
			}
		}
	}

	if err := runner.dockerClient.EnsureLocalImage(runner.ctx, "updater", imageRefs.Updater); err != nil {
		logger.WithError(err).Error("failed to ensure local image for updater")
		return
	}
	if err := runner.dockerClient.EnsureLocalImage(runner.ctx, "supervisor", imageRefs.Supervisor); err != nil {
		logger.WithError(err).Error("failed to ensure local image for supervisor")
		return
	}

	cfgBytes, err := json.Marshal(runner.cfg)
	if err != nil {
		logger.WithError(err).Error("cannot marshal config to json")
		return
	}
	cfgJson := string(cfgBytes)

	runner.updaterContainer, err = runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerUpdaterContainerName,
		Image: imageRefs.Updater,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "updater"},
		Env: map[string]string{
			config.EnvConfig: cfgJson,
		},
		Ports: map[string]string{
			runner.updaterPort: runner.updaterPort,
		},
		MaxLogSize:  runner.cfg.Log.MaxLogSize,
		MaxLogFiles: runner.cfg.Log.MaxLogFiles,
	})
	if err != nil {
		logger.WithError(err).Errorf("failed to start the updater")
		return
	}

	runner.supervisorContainer, err = runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerSupervisorContainerName,
		Image: imageRefs.Supervisor,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "supervisor"},
		Env: map[string]string{
			config.EnvConfig: cfgJson,
		},
		Volumes: map[string]string{
			"/var/run/docker.sock": "/var/run/docker.sock", // give access to host docker
			runner.cfg.FortaDir:    config.DefaultContainerFortaDirPath,
		},
		MaxLogSize:  runner.cfg.Log.MaxLogSize,
		MaxLogFiles: runner.cfg.Log.MaxLogFiles,
	})
	if err != nil {
		logger.WithError(err).Errorf("failed to start the supervisor")
		return
	}
}
