package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/feeds"
	"OpenZeppelin/fortify-node/services"
	"OpenZeppelin/fortify-node/services/registry"
	"OpenZeppelin/fortify-node/services/scanner"
	"OpenZeppelin/fortify-node/services/scanner/agentpool"
)

func loadKey() (*keystore.Key, error) {
	f, err := os.OpenFile("/passphrase", os.O_RDONLY, 400)
	if err != nil {
		return nil, err
	}

	pw, err := io.ReadAll(bufio.NewReader(f))
	if err != nil {
		return nil, err
	}
	passphrase := string(pw)

	files, err := ioutil.ReadDir("/.keys")
	if err != nil {
		return nil, err
	}

	if len(files) != 1 {
		return nil, errors.New("there must be only one key in key directory")
	}

	keyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", "/.keys", files[0].Name()))
	if err != nil {
		return nil, err
	}

	return keystore.DecryptKey(keyBytes, passphrase)
}

func initTxStream(ctx context.Context, cfg config.Config) (*scanner.TxStreamService, error) {
	url := cfg.Scanner.Ethereum.JsonRpcUrl
	startBlock := config.ParseBigInt(cfg.Scanner.StartBlock)
	endBlock := config.ParseBigInt(cfg.Scanner.EndBlock)
	chainID := config.ParseBigInt(cfg.Scanner.ChainID)

	if url == "" {
		return nil, fmt.Errorf("ethereum.jsonRpcUrl is required")
	}

	tracing := true
	if cfg.Scanner.DisableTracing {
		tracing = false
	}

	return scanner.NewTxStreamService(ctx, scanner.TxStreamServiceConfig{
		Url: url,
		BlockFeedConfig: feeds.BlockFeedConfig{
			Start:   startBlock,
			End:     endBlock,
			ChainID: chainID,
			Tracing: tracing,
		},
	})
}

func initTxAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool) (*scanner.TxAnalyzerService, error) {
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return scanner.NewTxAnalyzerService(ctx, scanner.TxAnalyzerServiceConfig{
		TxChannel:   stream.ReadOnlyTxStream(),
		AlertSender: as,
		AgentPool:   ap,
	})
}

func initBlockAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool) (*scanner.BlockAnalyzerService, error) {
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return scanner.NewBlockAnalyzerService(ctx, scanner.BlockAnalyzerServiceConfig{
		BlockChannel: stream.ReadOnlyBlockStream(),
		AlertSender:  as,
		AgentPool:    ap,
	})
}

func initAlertSender(ctx context.Context) (clients.AlertSender, error) {
	key, err := loadKey()
	if err != nil {
		return nil, err
	}
	qn := os.Getenv(config.EnvQueryNode)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvQueryNode)
	}
	return clients.NewAlertSender(ctx, clients.AlertSenderConfig{
		Key:           key,
		QueryNodeAddr: qn,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	messaging.Start("scanner", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	as, err := initAlertSender(ctx)
	if err != nil {
		return nil, err
	}
	txStream, err := initTxStream(ctx, cfg)
	if err != nil {
		return nil, err
	}

	agentPool := agentpool.NewAgentPool()
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream, agentPool)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream, agentPool)
	if err != nil {
		return nil, err
	}

	// Finally start the registry service so we know what agents we are running and receive updates.
	registryService := registry.New(cfg)

	return []services.Service{
		txStream,
		txAnalyzer,
		blockAnalyzer,
		scanner.NewTxLogger(ctx),
		registryService,
	}, nil
}

func main() {
	services.ContainerMain("scanner", initServices)
}
