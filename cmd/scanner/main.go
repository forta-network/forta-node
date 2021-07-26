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
	gethlog "github.com/ethereum/go-ethereum/log"

	"OpenZeppelin/fortify-node/clients"
	"OpenZeppelin/fortify-node/clients/messaging"
	"OpenZeppelin/fortify-node/config"
	"OpenZeppelin/fortify-node/ethereum"
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

func initTxStream(ctx context.Context, ethClient, traceClient ethereum.Client, cfg config.Config) (*scanner.TxStreamService, feeds.BlockFeed, error) {
	url := cfg.Scanner.Ethereum.JsonRpcUrl
	startBlock := config.ParseBigInt(cfg.Scanner.StartBlock)
	endBlock := config.ParseBigInt(cfg.Scanner.EndBlock)
	chainID := config.ParseBigInt(cfg.Scanner.ChainID)

	if url == "" {
		return nil, nil, fmt.Errorf("ethereum.jsonRpcUrl is required")
	}
	if cfg.Trace.Enabled && cfg.Trace.Ethereum.JsonRpcUrl == "" {
		return nil, nil, fmt.Errorf("trace requires a JsonRpcUrl if enabled")
	}

	blockFeed, err := feeds.NewBlockFeed(ctx, ethClient, traceClient, feeds.BlockFeedConfig{
		Start:   startBlock,
		End:     endBlock,
		ChainID: chainID,
		Tracing: cfg.Trace.Enabled,
	})
	if err != nil {
		return nil, nil, err
	}

	txStream, err := scanner.NewTxStreamService(ctx, ethClient, blockFeed, scanner.TxStreamServiceConfig{
		JsonRpcConfig:      cfg.Scanner.Ethereum,
		TraceJsonRpcConfig: cfg.Trace.Ethereum,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the tx stream service: %v", err)
	}

	return txStream, blockFeed, nil
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

func initRegistry(ctx context.Context, msgClient clients.MessageClient, cfg config.Config) (services.Service, error) {
	logFeed, err := feeds.NewLogFeed(ctx, cfg.Registry.Ethereum.WebsocketUrl, []string{cfg.Registry.ContractAddress})
	if err != nil {
		return nil, err
	}
	return registry.New(cfg, msgClient, logFeed), nil
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	as, err := initAlertSender(ctx)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethereum.NewStreamEthClient(ctx, cfg.Scanner.Ethereum.JsonRpcUrl)
	if err != nil {
		return nil, err
	}

	traceClient, err := ethereum.NewStreamEthClient(ctx, cfg.Trace.Ethereum.JsonRpcUrl)
	if err != nil {
		return nil, err
	}

	txStream, blockFeed, err := initTxStream(ctx, ethClient, traceClient, cfg)
	if err != nil {
		return nil, err
	}

	registryService, err := initRegistry(ctx, msgClient, cfg)
	if err != nil {
		return nil, err
	}

	agentPool := agentpool.NewAgentPool(msgClient)
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream, agentPool)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream, agentPool)
	if err != nil {
		return nil, err
	}

	// Start the main block feed so all transaction feeds can start consuming.
	blockFeed.Start()

	return []services.Service{
		txStream,
		txAnalyzer,
		blockAnalyzer,
		scanner.NewTxLogger(ctx),
		registryService,
	}, nil
}

func main() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)

	services.ContainerMain("scanner", initServices)
}
