package runner

import (
	"context"
	"strconv"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/store"
	"github.com/forta-protocol/forta-node/utils"
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
		logger := log.WithField("supervisor", latestRefs.Supervisor).WithField("updater", latestRefs.Updater)
		if latestRefs.Release != nil {
			logger = logger.WithField("commit", latestRefs.Release.Release.Commit)
		}
		logger.Info("received new node image reference")
		runner.replaceContainers(logger, latestRefs)
	}
}

func (runner *Runner) replaceContainers(logger *log.Entry, imageRefs store.ImageRefs) {
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

	for _, image := range []struct {
		Name string
		Ref  *string
	}{
		{
			Name: "updater",
			Ref:  &imageRefs.Updater,
		},
		{
			Name: "supervisor",
			Ref:  &imageRefs.Supervisor,
		},
	} {
		logger := log.WithField("ref", *image.Ref).WithField("name", image.Name)

		if !runner.cfg.Development {
			ref, err := utils.ValidateDiscoImageRef(runner.cfg.Registry.ContainerRegistry, *image.Ref)
			if err != nil {
				logger.WithError(err).Warn("not a disco ref - skipping pull")
				continue
			}
			*image.Ref = ref
		}
		// replace ref to include host in ref
		if err := runner.dockerClient.EnsureLocalImage(runner.ctx, image.Name, *image.Ref); err != nil {
			logger.WithError(err).Error("failed to ensure local image")
			return
		}
	}

	var err error

	runner.updaterContainer, err = runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerUpdaterContainerName,
		Image: imageRefs.Updater,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "updater"},
		Env: map[string]string{
			config.EnvDevelopment: strconv.FormatBool(runner.cfg.Development),
			config.EnvNoUpdate:    strconv.FormatBool(runner.cfg.NoUpdate),
		},
		Volumes: map[string]string{
			runner.cfg.FortaDir: config.DefaultContainerFortaDirPath,
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
	if err := runner.dockerClient.WaitContainerStart(runner.ctx, runner.updaterContainer.ID); err != nil {
		logger.WithError(err).Error("error while waiting for updater start")
		return
	}

	runner.supervisorContainer, err = runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerSupervisorContainerName,
		Image: imageRefs.Supervisor,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "supervisor"},
		Env: map[string]string{
			// supervisor needs to know and mount the forta dir on the host os
			config.EnvHostFortaDir: runner.cfg.FortaDir,
		},
		Volumes: map[string]string{
			// give access to host docker
			"/var/run/docker.sock": "/var/run/docker.sock",
			runner.cfg.FortaDir:    config.DefaultContainerFortaDirPath,
		},
		Files: map[string][]byte{
			"passphrase": []byte(runner.cfg.Passphrase),
		},
		MaxLogSize:  runner.cfg.Log.MaxLogSize,
		MaxLogFiles: runner.cfg.Log.MaxLogFiles,
	})
	if err != nil {
		logger.WithError(err).Errorf("failed to start the supervisor")
		return
	}
	if err := runner.dockerClient.WaitContainerStart(runner.ctx, runner.supervisorContainer.ID); err != nil {
		logger.WithError(err).Error("error while waiting for supervisor start")
		return
	}
}
