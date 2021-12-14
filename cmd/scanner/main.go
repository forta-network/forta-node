package scanner

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/feeds"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/registry"
	"github.com/forta-protocol/forta-node/services/scanner"
	"github.com/forta-protocol/forta-node/services/scanner/agentpool"
)

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

	var rateLimit *time.Ticker
	if cfg.Scanner.BlockRateLimit > 0 {
		rateLimit = time.NewTicker(time.Duration(cfg.Scanner.BlockRateLimit) * time.Millisecond)
	}

	blockFeed, err := feeds.NewBlockFeed(ctx, ethClient, traceClient, feeds.BlockFeedConfig{
		Start:     startBlock,
		End:       endBlock,
		ChainID:   chainID,
		Tracing:   cfg.Trace.Enabled,
		RateLimit: rateLimit,
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

func initTxAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool, msgClient clients.MessageClient) (*scanner.TxAnalyzerService, error) {
	qn := os.Getenv(config.EnvPublisherHost)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvPublisherHost)
	}
	return scanner.NewTxAnalyzerService(ctx, scanner.TxAnalyzerServiceConfig{
		TxChannel:   stream.ReadOnlyTxStream(),
		AlertSender: as,
		AgentPool:   ap,
		MsgClient:   msgClient,
	})
}

func initBlockAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool, msgClient clients.MessageClient) (*scanner.BlockAnalyzerService, error) {
	qn := os.Getenv(config.EnvPublisherHost)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvPublisherHost)
	}
	return scanner.NewBlockAnalyzerService(ctx, scanner.BlockAnalyzerServiceConfig{
		BlockChannel: stream.ReadOnlyBlockStream(),
		AlertSender:  as,
		AgentPool:    ap,
		MsgClient:    msgClient,
	})
}

func initAlertSender(ctx context.Context, key *keystore.Key) (clients.AlertSender, error) {
	qn := os.Getenv(config.EnvPublisherHost)
	if qn == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvPublisherHost)
	}
	return clients.NewAlertSender(ctx, clients.AlertSenderConfig{
		Key:           key,
		QueryNodeAddr: qn,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	cfg.LocalAgentsPath = config.DefaultContainerLocalAgentsFilePath

	natsHost := os.Getenv(config.EnvNatsHost)
	if natsHost == "" {
		return nil, fmt.Errorf("%s is a required env var", config.EnvNatsHost)
	}
	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", natsHost, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	as, err := initAlertSender(ctx, key)
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

	registryService := registry.New(cfg, key.Address, msgClient)
	agentPool := agentpool.NewAgentPool(cfg.Scanner, msgClient)
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}

	// Start the main block feed so all transaction feeds can start consuming.
	if !cfg.Scanner.DisableAutostart {
		blockFeed.Start()
	}

	svcs := []services.Service{
		txStream,
		txAnalyzer,
		blockAnalyzer,
		scanner.NewScannerAPI(ctx, blockFeed),
		scanner.NewTxLogger(ctx),
	}

	// for performance tests, this flag avoids using registry service
	if !cfg.Registry.Disabled {
		svcs = append(svcs, registryService)
	}

	return svcs, nil
}

func Run() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)
	services.ContainerMain("scanner", initServices)
}
