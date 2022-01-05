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

	currentUpdaterImg    string
	currentSupervisorImg string

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
	runner.removeContainer(runner.updaterContainer)
	runner.removeContainer(runner.supervisorContainer)
	return nil
}

func (runner *Runner) removeContainer(container *clients.DockerContainer) error {
	if container != nil {
		logger := log.WithField("container", container.ID).WithField("name", container.Name)
		if err := runner.dockerClient.InterruptContainer(context.Background(), container.ID); err != nil {
			logger.WithError(err).Error("error stopping container")
		} else {
			logger.Info("interrupted")
		}
		if err := runner.dockerClient.WaitContainerExit(context.Background(), container.ID); err != nil {
			logger.WithError(err).Panic("error while waiting for container exit")
			return err
		}
		if err := runner.dockerClient.Prune(runner.ctx); err != nil {
			logger.WithError(err).Panic("error while pruning after stopping old containers")
			return err
		}
		if err := runner.dockerClient.WaitContainerPrune(runner.ctx, container.ID); err != nil {
			logger.WithError(err).Panic("error while waiting for old container prune")
			return err
		}
	}
	return nil
}

func (runner *Runner) receive() {
	for latestRefs := range runner.imgStore.Latest() {
		logger := log.WithField("supervisor", latestRefs.Supervisor).WithField("updater", latestRefs.Updater)
		if latestRefs.ReleaseInfo != nil {
			logger = logger.WithFields(log.Fields{
				"commit":      latestRefs.ReleaseInfo.Manifest.Release.Commit,
				"releaseInfo": latestRefs.ReleaseInfo.String(),
			})
		}
		logger.Info("detected new images")
		if latestRefs.Updater != runner.currentUpdaterImg {
			if err := runner.replaceUpdater(logger, latestRefs); err != nil {
				logger.WithError(err).Panic("error replacing updater")
			} else {
				runner.currentUpdaterImg = latestRefs.Updater
			}
		}

		// in local mode, run supervisor even if release info is nil.
		// in production mode, it doesn't make sense to run the supervisor for the first time yet
		// if updater hasn't yet received the latest release.
		validReleaseInfo := latestRefs.ReleaseInfo != nil || config.UseDockerImages == "local"
		if latestRefs.Supervisor != runner.currentSupervisorImg && validReleaseInfo {
			if err := runner.replaceSupervisor(logger, latestRefs); err != nil {
				logger.WithError(err).Panic("error replacing supervisor")
			} else {
				runner.currentSupervisorImg = latestRefs.Supervisor
			}
		}
	}
}

func (runner *Runner) ensureImage(logger *log.Entry, name string, imageRef string) (string, error) {
	logger = logger.WithField("ref", imageRef).WithField("name", name)

	// to make things easier, don't require image ref validation in dev mode
	if !runner.cfg.Development {
		fixedRef, err := utils.ValidateDiscoImageRef(runner.cfg.Registry.ContainerRegistry, imageRef)
		if err != nil {
			logger.WithError(err).Warn("not a disco ref - skipping pull")
		} else {
			imageRef = fixedRef // important
		}
	}

	if err := runner.dockerClient.EnsureLocalImage(runner.ctx, name, imageRef); err != nil {
		logger.WithError(err).Warn("failed to ensure local image")
		return "", err
	}

	return imageRef, nil
}

func (runner *Runner) replaceUpdater(logger *log.Entry, imageRefs store.ImageRefs) error {
	logger.Info("replacing updater")
	err := runner.removeContainer(runner.updaterContainer)
	if err != nil {
		return err
	}
	return runner.startUpdater(logger, imageRefs)
}

func (runner *Runner) replaceSupervisor(logger *log.Entry, imageRefs store.ImageRefs) error {
	logger.Info("replacing supervisor")
	err := runner.removeContainer(runner.supervisorContainer)
	if err != nil {
		return err
	}
	return runner.startSupervisor(logger, imageRefs)
}

func (runner *Runner) startUpdater(logger *log.Entry, latestRefs store.ImageRefs) (err error) {
	updaterRef := latestRefs.Updater
	updaterRef, err = runner.ensureImage(logger, "updater", updaterRef)
	if err != nil {
		return err
	}

	uc, err := runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerUpdaterContainerName,
		Image: updaterRef,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "updater"},
		Env: map[string]string{
			config.EnvDevelopment: strconv.FormatBool(runner.cfg.Development),
			config.EnvNoUpdate:    strconv.FormatBool(runner.cfg.NoUpdate),
			config.EnvReleaseInfo: latestRefs.ReleaseInfo.String(),
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
		return err
	}
	runner.updaterContainer = uc

	if err := runner.dockerClient.WaitContainerStart(runner.ctx, runner.updaterContainer.ID); err != nil {
		logger.WithError(err).Error("error while waiting for updater start")
		return err
	}
	return nil
}

func (runner *Runner) startSupervisor(logger *log.Entry, latestRefs store.ImageRefs) (err error) {
	supervisorRef := latestRefs.Supervisor
	supervisorRef, err = runner.ensureImage(logger, "supervisor", supervisorRef)
	if err != nil {
		return err
	}
	sc, err := runner.dockerClient.StartContainer(runner.ctx, clients.DockerContainerConfig{
		Name:  config.DockerSupervisorContainerName,
		Image: supervisorRef,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "supervisor"},
		Env: map[string]string{
			// supervisor needs to know and mount the forta dir on the host os
			config.EnvHostFortaDir: runner.cfg.FortaDir,
			config.EnvReleaseInfo:  latestRefs.ReleaseInfo.String(),
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
		return err
	}
	runner.supervisorContainer = sc

	if err := runner.dockerClient.WaitContainerStart(runner.ctx, runner.supervisorContainer.ID); err != nil {
		logger.WithError(err).Error("error while waiting for supervisor start")
		return err
	}
	return nil
}
