package services

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/forta-protocol/forta-node/ens"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"github.com/forta-protocol/forta-node/config"
)

// Service is a service abstraction.
type Service interface {
	Start() error
	Stop() error
	Name() string
}

var sigc chan os.Signal

var execIDKey = struct{}{}

func ExecID(ctx context.Context) string {
	execID := ctx.Value(execIDKey)
	if execID == nil {
		panic("cannot get exec ID")
	}
	return execID.(string)
}

func initExecID(ctx context.Context) context.Context {
	execID, err := uuid.NewUUID()
	if err != nil {
		panic(err)
	}
	return context.WithValue(ctx, execIDKey, execID.String())
}

func setContracts(cfg *config.Config) error {
	contracts, err := ens.ResolveFortaContracts(cfg.ENSConfig.JsonRpc.Url, cfg.ENSConfig.ContractAddress)
	if err != nil {
		return err
	}
	if cfg.Registry.ContractAddress == "" {
		cfg.Registry.ContractAddress = contracts.Dispatch
	}
	cfg.ScannerVersionContractAddress = contracts.ScannerVersion
	cfg.AgentRegistryContractAddress = contracts.Agent
	return nil
}

func ContainerMain(name string, getServices func(ctx context.Context, cfg config.Config) ([]Service, error)) {
	logger := log.WithField("container", name)

	cfg, err := config.GetConfigForContainer()
	if err != nil {
		logger.WithError(err).Error("could not get config")
		return
	}

	if err := setContracts(&cfg); err != nil {
		logger.WithError(err).Error("could not initialize contract addresses using config")
		return
	}

	lvl, err := log.ParseLevel(cfg.Log.Level)
	if err != nil {
		logger.WithError(err).Error("could not initialize log level")
		return
	}
	log.SetLevel(lvl)
	log.SetFormatter(&log.JSONFormatter{})
	logger.Info("starting")
	defer logger.Info("exiting")

	ctx, cancel := InitMainContext()
	defer cancel()

	serviceList, err := getServices(ctx, cfg)
	if err != nil {
		logger.WithError(err).Error("could not initialize services")
		return
	}

	if err := StartServices(ctx, cancel, logger, serviceList); err != nil {
		logger.WithError(err).Error("failed to start services")
	}
}

func InitMainContext() (context.Context, context.CancelFunc) {
	execIDCtx := initExecID(context.Background())
	ctx, cancel := context.WithCancel(execIDCtx)
	if sigc == nil {
		sigc = make(chan os.Signal, 1)
	}
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	go func() {
		sig := <-sigc
		log.Infof("received signal: %s", sig.String())
		cancel()
	}()
	return ctx, cancel
}

// StartServices kicks off all services.
func StartServices(ctx context.Context, cancelMainCtx context.CancelFunc, logger *log.Entry, services []Service) error {
	// each service should be able to start successfully within reasonable time
	for _, service := range services {
		serviceStartedCtx, serviceStarted := context.WithCancel(context.Background())
		defer serviceStarted()

		logger := logger.WithField("service", service.Name())

		go func() {
			if err := service.Start(); err != nil {
				logger.WithError(err).Error("failed to start service")
				cancelMainCtx()
				return
			}
			serviceStarted()
		}()

		select {
		case <-time.After(time.Minute):
			logger.Error("took too long to start service")
			cancelMainCtx()
			break
		case <-serviceStartedCtx.Done():
			// ok - do nothing
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	<-ctx.Done()
	logger.WithError(ctx.Err()).Info("context is done")

	// stop all services
	for _, service := range services {
		err := service.Stop()
		logger.WithError(err).WithField("service", service.Name()).Info("stopped")
	}

	return nil
}
