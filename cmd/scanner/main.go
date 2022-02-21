package scanner

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/health"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/ethereum"
	"github.com/forta-protocol/forta-node/feeds"
	"github.com/forta-protocol/forta-node/security"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/registry"
	"github.com/forta-protocol/forta-node/services/scanner"
	"github.com/forta-protocol/forta-node/services/scanner/agentpool"
	"github.com/forta-protocol/forta-node/utils"
)

func initTxStream(ctx context.Context, ethClient, traceClient ethereum.Client, cfg config.Config) (*scanner.TxStreamService, feeds.BlockFeed, error) {
	url := cfg.Scan.JsonRpc.Url
	chainID := config.ParseBigInt(cfg.ChainID)

	if url == "" {
		return nil, nil, fmt.Errorf("scan.jsonRpc.url is required")
	}
	if cfg.Trace.Enabled && cfg.Trace.JsonRpc.Url == "" {
		return nil, nil, fmt.Errorf("trace requires a jsonRpc URL if enabled")
	}

	var rateLimit *time.Ticker
	if cfg.Scan.BlockRateLimit > 0 {
		rateLimit = time.NewTicker(time.Duration(cfg.Scan.BlockRateLimit) * time.Millisecond)
	}

	var maxAge time.Duration
	if cfg.Scan.BlockMaxAgeSeconds > 0 {
		maxAge = time.Duration(cfg.Scan.BlockMaxAgeSeconds) * time.Second
	}
	blockFeed, err := feeds.NewBlockFeed(ctx, ethClient, traceClient, feeds.BlockFeedConfig{
		ChainID:             chainID,
		Tracing:             cfg.Trace.Enabled,
		RateLimit:           rateLimit,
		SkipBlocksOlderThan: &maxAge,
	})
	if err != nil {
		return nil, nil, err
	}

	txStream, err := scanner.NewTxStreamService(ctx, ethClient, blockFeed, scanner.TxStreamServiceConfig{
		JsonRpcConfig:       cfg.Scan.JsonRpc,
		TraceJsonRpcConfig:  cfg.Trace.JsonRpc,
		SkipBlocksOlderThan: &maxAge,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the tx stream service: %v", err)
	}

	return txStream, blockFeed, nil
}

func initTxAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool, msgClient clients.MessageClient) (*scanner.TxAnalyzerService, error) {
	return scanner.NewTxAnalyzerService(ctx, scanner.TxAnalyzerServiceConfig{
		TxChannel:   stream.ReadOnlyTxStream(),
		AlertSender: as,
		AgentPool:   ap,
		MsgClient:   msgClient,
	})
}

func initBlockAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.TxStreamService, ap *agentpool.AgentPool, msgClient clients.MessageClient) (*scanner.BlockAnalyzerService, error) {
	return scanner.NewBlockAnalyzerService(ctx, scanner.BlockAnalyzerServiceConfig{
		BlockChannel: stream.ReadOnlyBlockStream(),
		AlertSender:  as,
		AgentPool:    ap,
		MsgClient:    msgClient,
	})
}

func initAlertSender(ctx context.Context, key *keystore.Key) (clients.AlertSender, error) {
	return clients.NewAlertSender(ctx, clients.AlertSenderConfig{
		Key:               key,
		PublisherNodeAddr: config.DockerPublisherContainerName,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	cfg.LocalAgentsPath = config.DefaultContainerLocalAgentsFilePath

	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.Trace.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Trace.JsonRpc.Url)

	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	as, err := initAlertSender(ctx, key)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethereum.NewStreamEthClient(ctx, "chain", cfg.Scan.JsonRpc.Url)
	if err != nil {
		return nil, err
	}

	traceClient, err := ethereum.NewStreamEthClient(ctx, "trace", cfg.Trace.JsonRpc.Url)
	if err != nil {
		return nil, err
	}

	txStream, blockFeed, err := initTxStream(ctx, ethClient, traceClient, cfg)
	if err != nil {
		return nil, err
	}

	registryService := registry.New(cfg, key.Address, msgClient, ethClient)
	agentPool := agentpool.NewAgentPool(ctx, cfg.Scan, msgClient)
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}

	// Start the main block feed so all transaction feeds can start consuming.
	if !cfg.Scan.DisableAutostart {
		blockFeed.Start()
	}

	svcs := []services.Service{
		health.NewService(ctx, health.CheckerFrom(
			summarizeReports,
			ethClient, traceClient, blockFeed, txStream, txAnalyzer, blockAnalyzer, agentPool, registryService,
		)),
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

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	agentsTotal, ok := reports.NameContains("agent-pool.agents.total")
	var (
		gotCount bool
		count    int
	)
	var err error
	if ok {
		count, err = strconv.Atoi(agentsTotal.Details)
		if err != nil {
			summary.Addf("running %d agents", count)
			gotCount = true
		}
	}

	checkedTime, ok := reports.NameContains("registry.event.checked.time")
	var checkedT *time.Time
	if ok && len(checkedTime.Details) > 0 {
		checkedT, _ = checkedTime.Time()
	}
	var checkedE string
	checkedErr, ok := reports.NameContains("registry.event.checked.time")
	if ok {
		checkedE = checkedErr.Details
	}
	if gotCount && count == 0 {
		if len(checkedE) > 0 {
			summary.Addf("because failed to check the list with error '%s'", checkedE)
			summary.Status(health.StatusFailing)
		}
		if checkedTime != nil {
			checkDelay := time.Since(*checkedT)
			if checkDelay > time.Minute*10 {
				summary.Addf("and delayed for %d minutes", int64(checkDelay.Minutes()))
				summary.Status(health.StatusFailing)
			}
		}
	}
	summary.Punc(".")

	lastBlock, ok := reports.NameContains("block-feed.last-block")
	if ok && len(lastBlock.Details) > 0 {
		summary.Addf("at block %s.", lastBlock.Details)
	}

	blockByNumberErr, ok := reports.NameContains("chain-json-rpc-client.request.block-by-number.error")
	if ok && len(blockByNumberErr.Details) > 0 {
		summary.Addf("failing to get block with error '%s'", blockByNumberErr.Details)
		summary.Status(health.StatusFailing)
	}
	blockByNumberTime, ok := reports.NameContains("chain-json-rpc-client.request.block-by-number.error")
	if ok && len(blockByNumberTime.Details) > 0 {
		t, ok := blockByNumberTime.Time()
		if ok {
			checkDelay := time.Since(*t)
			if checkDelay > time.Minute*5 {
				summary.Punc(",")
				summary.Addf("lagging to get block by %d minutes.", checkDelay.Minutes())
			}
		}
	}
	summary.Punc(".")

	getTxReceiptErr, ok := reports.NameContains("chain-json-rpc-client.request.get-transaction-receipt")
	if ok && len(getTxReceiptErr.Details) > 0 {
		summary.Addf("failing to get transaction receipt with error '%s', this can slow down block processing.", getTxReceiptErr.Details)
	}

	traceBlockErr, ok := reports.NameContains("trace-json-rpc-client.request.trace-block.error")
	if ok && len(traceBlockErr.Details) > 0 {
		summary.Addf("trace api (trace_block) is failing with error '%s'.", traceBlockErr.Details)
	}

	return summary.Finish()
}

func Run() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)
	services.ContainerMain("scanner", initServices)
}
