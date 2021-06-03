package main

import (
	"context"
	"flag"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	log "github.com/sirupsen/logrus"

	"github.com/OpenZeppelin/fortify-node/config"
	"github.com/OpenZeppelin/fortify-node/services"
	"github.com/OpenZeppelin/fortify-node/services/containers"
)

func initServices(cfg config.Config, passphrase string, ctx context.Context) ([]services.Service, error) {
	svc, err := containers.NewTxNodeService(ctx, containers.TxNodeServiceConfig{
		Config:     cfg,
		Passphrase: passphrase,
	})
	if err != nil {
		return nil, err
	}
	return []services.Service{
		svc,
	}, nil
}

func initKeyFile(passphrase string) error {
	keyPath, err := config.GetKeyStorePath()
	if err != nil {
		return err
	}
	f, err := os.OpenFile(keyPath, os.O_RDONLY, 400)
	defer func() {
		if f != nil {
			f.Close()
		}
	}()
	if os.IsNotExist(err) {
		ks := keystore.NewKeyStore(keyPath, keystore.StandardScryptN, keystore.StandardScryptP)
		acct, err := ks.NewAccount(passphrase)
		if err != nil {
			return err
		}
		log.Infof("Generated Address: %s", acct.Address.Hex())
	}
	return nil
}

func main() {

	ctx, cancel := services.InitMainContext()
	defer cancel()

	cfgFile := flag.String("config", "config.yml", "filename for configuration yaml")
	passphrase := flag.String("passphrase", "", "passphrase for configuration yaml")

	flag.Parse()

	if *passphrase == "" {
		log.Error("-passphrase is required")
		return
	}

	if err := initKeyFile(*passphrase); err != nil {
		log.Error("could not initialize key", err)
		return
	}

	cfg, err := config.GetConfig(*cfgFile)
	if err != nil {
		log.Error("could not read config file", err)
		return
	}
	if err := config.InitLogLevel(cfg); err != nil {
		log.Error("error initializing log level", err)
		return
	}

	log.Info("Starting Node")

	serviceList, err := initServices(cfg, *passphrase, ctx)
	if err != nil {
		log.Error("could not initialize services", err)
		return
	}

	if err := services.StartServices(ctx, serviceList); err != nil {
		log.Error("error running services", err)
	}

	log.Info("Stopping Node")
}
