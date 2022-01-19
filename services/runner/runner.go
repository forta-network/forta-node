package runner

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/alertapi"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/domain"
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
	nukeClient   clients.DockerClient

	// test clients for initial checks
	scanClient  EthereumClient
	traceClient EthereumClient
	alertClient clients.AlertAPIClient

	currentUpdaterImg    string
	currentSupervisorImg string

	updaterPort         string
	updaterContainer    *clients.DockerContainer
	supervisorContainer *clients.DockerContainer
}

// EthereumClient is useful for checking the JSON-RPC API.
type EthereumClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

// NewRunner creates a new runner.
func NewRunner(ctx context.Context, cfg config.Config,
	imgStore store.FortaImageStore, runnerDockerClient clients.DockerClient,
	globalDockerClient clients.DockerClient, updaterPort string,
) *Runner {
	return &Runner{
		ctx:          ctx,
		cfg:          cfg,
		imgStore:     imgStore,
		dockerClient: runnerDockerClient,
		nukeClient:   globalDockerClient,
		updaterPort:  updaterPort,
	}
}

// Start starts the service.
func (runner *Runner) Start() error {
	if err := runner.setupTestClients(); err != nil {
		return err
	}
	if err := runner.doStartUpCheck(); err != nil {
		return fmt.Errorf("start-up check failed: %v", err)
	}
	log.Info("start-up check successful")

	if err := runner.nukeClient.Nuke(context.Background()); err != nil {
		return fmt.Errorf("failed to nuke leftover containers at start: %v", err)
	}

	go runner.receive()
	return nil
}

// Name returns the name of the service.
func (runner *Runner) Name() string {
	return "runner"
}

// Stop stops the service
func (runner *Runner) Stop() error {
	if err := runner.nukeClient.Nuke(context.Background()); err != nil {
		return fmt.Errorf("failed to nuke containers before exiting: %v", err)
	}
	return nil
}

func (runner *Runner) setupTestClients() error {
	var err error
	scanClientUrl, err := runner.fixTestRpcUrl(runner.cfg.Scan.JsonRpc.Url)
	if err != nil {
		return fmt.Errorf("invalid scan url: %v", err)
	}
	runner.scanClient, err = ethclient.Dial(scanClientUrl)
	if err != nil {
		return err
	}
	traceClientUrl, err := runner.fixTestRpcUrl(runner.cfg.Trace.JsonRpc.Url)
	if err != nil {
		return fmt.Errorf("invalid trace url: %v", err)
	}
	runner.traceClient, err = ethclient.Dial(traceClientUrl)
	if err != nil {
		return err
	}
	runner.alertClient = alertapi.NewClient(runner.cfg.Publish.APIURL)
	return nil
}

func (runner *Runner) doStartUpCheck() error {
	// ensure that docker is available
	_, err := runner.dockerClient.GetContainers(runner.ctx)
	if err != nil {
		return fmt.Errorf("docker check failed (get containers): %v", err)
	}
	// ensure that the scan json-rpc api is reachable
	_, err = runner.scanClient.BlockNumber(runner.ctx)
	if err != nil {
		return fmt.Errorf("scan api check failed: %v", err)
	}
	// ensure that the trace json-rpc api is reachable
	_, err = runner.traceClient.BlockNumber(runner.ctx)
	if err != nil {
		return fmt.Errorf("trace api check failed: %v", err)
	}
	// ensure that the batch api is available for publishing to
	if err := runner.alertClient.PostBatch(&domain.AlertBatch{Ref: "test"}, ""); err != nil {
		return fmt.Errorf("batch api check failed: %v", err)
	}
	return nil
}

func (runner *Runner) fixTestRpcUrl(rawurl string) (string, error) {
	return strings.ReplaceAll(rawurl, "host.docker.internal", "localhost"), nil
}

func (runner *Runner) removeContainer(container *clients.DockerContainer) error {
	if container != nil {
		return runner.removeContainerWithProps(container.Name, container.ID)
	}
	return nil
}

func (runner *Runner) removeContainerWithProps(name, id string) error {
	logger := log.WithField("container", id).WithField("name", name)
	if err := runner.dockerClient.InterruptContainer(context.Background(), id); err != nil {
		logger.WithError(err).Error("error stopping container")
	} else {
		logger.Info("interrupted")
	}
	if err := runner.dockerClient.WaitContainerExit(context.Background(), id); err != nil {
		logger.WithError(err).Panic("error while waiting for container exit")
	}
	if err := runner.dockerClient.Prune(runner.ctx); err != nil {
		logger.WithError(err).Panic("error while pruning after stopping old containers")
	}
	if err := runner.dockerClient.WaitContainerPrune(runner.ctx, id); err != nil {
		logger.WithError(err).Panic("error while waiting for old container prune")
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
