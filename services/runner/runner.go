package runner

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/docker"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/components/prometheus"
	"github.com/forta-network/forta-node/store"
	log "github.com/sirupsen/logrus"
)

// Errprs
var (
	ErrBadProxyAPI = fmt.Errorf("proxy api must be specified as http(s) when scan api is websocket")
)

// Runner receives and starts the latest updater and supervisor.
type Runner struct {
	ctx          context.Context
	cfg          config.Config
	imgStore     store.FortaImageStore
	dockerClient clients.DockerClient
	globalClient clients.DockerClient

	updaterContainer     *docker.Container
	supervisorContainer  *docker.Container
	currentUpdaterImg    string
	currentSupervisorImg string
	containerMu          sync.RWMutex // protects above refs and containers

	healthClient health.HealthClient
}

// EthereumClient is useful for checking the JSON-RPC API.
type EthereumClient interface {
	BlockNumber(ctx context.Context) (uint64, error)
}

// NewRunner creates a new runner.
func NewRunner(ctx context.Context, cfg config.Config,
	imgStore store.FortaImageStore, runnerDockerClient clients.DockerClient,
	globalDockerClient clients.DockerClient,
) *Runner {
	return &Runner{
		ctx:          ctx,
		cfg:          cfg,
		imgStore:     imgStore,
		dockerClient: runnerDockerClient,
		globalClient: globalDockerClient,
		healthClient: health.NewClient(),
	}
}

// Start starts the service.
func (runner *Runner) Start() error {
	if err := runner.doStartUpCheck(); err != nil {
		return fmt.Errorf("start-up check failed: %v", err)
	}
	log.Info("start-up check successful")

	if err := runner.globalClient.Nuke(context.Background()); err != nil {
		return fmt.Errorf("failed to nuke leftover containers at start: %v", err)
	}

	health.StartServer(runner.ctx, "", healthutils.DefaultHealthServerErrHandler, runner.CheckServiceHealth)

	if runner.cfg.AutoUpdate.Disable {
		runner.startEmbeddedSupervisor()
	} else {
		runner.startEmbeddedUpdater()
		go runner.keepContainersUpToDate()
	}

	go runner.keepContainersAlive()

	prometheus.StartCollector(runner, runner.cfg.PrometheusConfig.Port)

	return nil
}

// Name returns the name of the service.
func (runner *Runner) Name() string {
	return "runner"
}

// Stop stops the service
func (runner *Runner) Stop() error {
	runner.containerMu.RLock()
	defer runner.containerMu.RUnlock()

	if runner.updaterContainer != nil {
		runner.dockerClient.InterruptContainer(context.Background(), runner.updaterContainer.ID)

	}
	if runner.supervisorContainer != nil {
		runner.dockerClient.InterruptContainer(context.Background(), runner.supervisorContainer.ID)
	}
	return nil
}

func (runner *Runner) doStartUpCheck() error {
	// ensure that docker is available
	_, err := runner.dockerClient.GetContainers(runner.ctx)
	if err != nil {
		return fmt.Errorf("docker check failed (get containers): %v", err)
	}
	// ensure that the scan json-rpc api is reachable
	err = ethereum.TestAPI(runner.ctx, runner.fixTestRpcUrl(runner.cfg.Scan.JsonRpc.Url))
	if err != nil {
		return fmt.Errorf("scan api check failed: %v", err)
	}

	if err := CheckProxyAgainstScan(runner.cfg.Scan.JsonRpc.Url, runner.cfg.JsonRpcProxy.JsonRpc.Url); err != nil {
		return err
	}

	if runner.cfg.Trace.Enabled && !runner.cfg.LocalModeConfig.Enable {
		// ensure that the trace json-rpc api is reachable
		err = ethereum.TestAPI(runner.ctx, runner.fixTestRpcUrl(runner.cfg.Trace.JsonRpc.Url))
		if err != nil {
			return fmt.Errorf("trace api check failed: %v", err)
		}
	}
	return nil
}

func (runner *Runner) fixTestRpcUrl(rawurl string) string {
	return strings.ReplaceAll(rawurl, "host.docker.internal", "localhost")
}

func (runner *Runner) removeContainer(container *docker.Container) error {
	if container != nil {
		return runner.removeContainerWithProps(container.Name, container.ID)
	}
	return nil
}

func (runner *Runner) removeContainerWithProps(name, id string) error {
	logger := log.WithField("container", id).WithField("name", name)
	if err := runner.dockerClient.TerminateContainer(context.Background(), id); err != nil {
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

func (runner *Runner) startEmbeddedUpdater() {
	builtInRefs := runner.imgStore.EmbeddedImageRefs()
	logger := log.WithField("supervisor", builtInRefs.Supervisor).WithField("updater", builtInRefs.Updater)

	if err := runner.replaceUpdater(logger, builtInRefs); err != nil {
		logger.WithError(err).Panic("error replacing updater")
	} else {
		runner.currentUpdaterImg = builtInRefs.Updater
	}
}

func (runner *Runner) startEmbeddedSupervisor() {
	builtInRefs := runner.imgStore.EmbeddedImageRefs()
	logger := log.WithField("supervisor", builtInRefs.Supervisor).WithField("updater", builtInRefs.Updater)

	if err := runner.replaceSupervisor(logger, builtInRefs); err != nil {
		logger.WithError(err).Panic("error replacing supervisor")
	} else {
		runner.currentSupervisorImg = builtInRefs.Supervisor
	}
}

func (runner *Runner) keepContainersUpToDate() {
	defer func() {
		if r := recover(); r != nil {
			runner.Stop()
			panic(r)
		}
	}()

	for latestRefs := range runner.imgStore.Latest() {
		runner.updateContainers(latestRefs)
	}
}

func (runner *Runner) updateContainers(latestRefs store.ImageRefs) {
	runner.containerMu.Lock()
	defer runner.containerMu.Unlock()

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
	} else {
		log.Debug("same image - not replacing updater")
	}

	if latestRefs.Supervisor != runner.currentSupervisorImg {
		if err := runner.replaceSupervisor(logger, latestRefs); err != nil {
			logger.WithError(err).Panic("error replacing supervisor")
		} else {
			runner.currentSupervisorImg = latestRefs.Supervisor
		}
	} else {
		log.Debug("same image - not replacing supervisor")
	}
}

func (runner *Runner) ensureImage(logger *log.Entry, name string, imageRef string) (string, error) {
	logger = logger.WithField("ref", imageRef).WithField("name", name)

	// to make things easier, don't require image ref validation in dev mode
	if !runner.cfg.Development {
		fixedRef, err := utils.ValidateDiscoImageRef(runner.cfg.Registry.ContainerRegistry, imageRef)
		if err != nil {
			logger.WithError(err).WithField("imageRef", imageRef).Warn("not a disco ref")
		} else {
			imageRef = fixedRef // important
		}
	}

	ticker := time.NewTicker(time.Minute)

	for {
		err := runner.dockerClient.EnsureLocalImage(runner.ctx, name, imageRef)
		if err != nil {
			logger.WithError(err).Warn("failed to ensure local image - retrying")
		} else {
			break
		}
		select {
		case <-ticker.C:
			// continue
		case <-runner.ctx.Done():
			// returning underlying err, because it's != nil
			return "", err
		}
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

	uc, err := runner.dockerClient.StartContainer(runner.ctx, docker.ContainerConfig{
		Name:  config.DockerUpdaterContainerName,
		Image: updaterRef,
		Cmd:   []string{config.DefaultFortaNodeBinaryPath, "updater"},
		Env: map[string]string{
			config.EnvDevelopment: strconv.FormatBool(runner.cfg.Development),
			config.EnvReleaseInfo: latestRefs.ReleaseInfo.String(),
		},
		Volumes: map[string]string{
			runner.cfg.FortaDir: config.DefaultContainerFortaDirPath,
		},
		Ports: map[string]string{
			config.DefaultContainerPort: config.DefaultContainerPort,
			"":                          config.DefaultHealthPort, // random host port
		},
		DialHost:    true,
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
	sc, err := runner.dockerClient.StartContainer(runner.ctx, docker.ContainerConfig{
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
		Ports: map[string]string{
			"": config.DefaultHealthPort, // random host port
		},
		Files: map[string][]byte{
			"passphrase": []byte(runner.cfg.Passphrase),
		},
		DialHost:    true,
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

func (runner *Runner) keepContainersAlive() {
	ticker := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ticker.C:
			if err := runner.doKeepContainersAlive(); err != nil {
				log.WithError(err).Error("failed while keeping containers alive")
			}

		case <-runner.ctx.Done():
			return
		}
	}
}

func (runner *Runner) doKeepContainersAlive() error {
	runner.containerMu.Lock()
	defer runner.containerMu.Unlock()

	if runner.supervisorContainer != nil {
		container, err := runner.dockerClient.GetContainerByID(runner.ctx, runner.supervisorContainer.ID)
		if err == nil && container.State == "exited" {
			containerDetails, err := runner.dockerClient.InspectContainer(runner.ctx, container.ID)
			if err != nil {
				return err
			}
			if containerDetails.State.ExitCode == services.ExitCodeTriggered {
				log.WithField("name", runner.supervisorContainer.Name).Info("detected internal exit trigger - exiting")
				services.TriggerExit(0)
				return nil
			}
			runner.dockerClient.StartContainer(runner.ctx, runner.supervisorContainer.Config)
		}
	}

	// only keep updater up if auto-update is enabled
	if runner.updaterContainer != nil && !runner.cfg.AutoUpdate.Disable {
		container, err := runner.dockerClient.GetContainerByID(runner.ctx, runner.updaterContainer.ID)
		if err == nil && container.State == "exited" {
			runner.dockerClient.StartContainer(runner.ctx, runner.updaterContainer.Config)
		}
	}

	return nil
}
