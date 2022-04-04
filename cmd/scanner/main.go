package scanner

import (
	"context"
	"fmt"
	"github.com/forta-protocol/forta-node/services/publisher"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-protocol/forta-core-go/clients/health"
	"github.com/forta-protocol/forta-core-go/ethereum"
	"github.com/forta-protocol/forta-core-go/feeds"
	"github.com/forta-protocol/forta-core-go/security"
	"github.com/forta-protocol/forta-core-go/utils"
	"github.com/forta-protocol/forta-node/clients"
	"github.com/forta-protocol/forta-node/clients/messaging"
	"github.com/forta-protocol/forta-node/config"
	"github.com/forta-protocol/forta-node/healthutils"
	"github.com/forta-protocol/forta-node/services"
	"github.com/forta-protocol/forta-node/services/registry"
	"github.com/forta-protocol/forta-node/services/scanner"
	"github.com/forta-protocol/forta-node/services/scanner/agentpool"
)

func initTxStream(ctx context.Context, ethClient, traceClient ethereum.Client, cfg config.Config) (*scanner.TxStreamService, feeds.BlockFeed, error) {
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)
	cfg.Registry.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.APIURL)
	cfg.Registry.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.GatewayURL)

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
		Offset:              config.GetBlockOffset(cfg.ChainID),
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

func initAlertSender(ctx context.Context, key *keystore.Key, pubClient clients.PublishClient) (clients.AlertSender, error) {
	return clients.NewLocalAlertSender(ctx, pubClient, clients.AlertSenderConfig{
		Key:               key,
		PublisherNodeAddr: config.DockerPublisherContainerName,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	cfg.LocalAgentsPath = config.DefaultContainerLocalAgentsFilePath

	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.Trace.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Trace.JsonRpc.Url)
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)
	cfg.Registry.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.APIURL)
	cfg.Registry.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.GatewayURL)

	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	cfg.Publish.APIURL = utils.ConvertToDockerHostURL(cfg.Publish.APIURL)
	cfg.Publish.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Publish.IPFS.APIURL)
	cfg.Publish.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Publish.IPFS.GatewayURL)

	publisherSvc, err := publisher.NewPublisher(ctx, cfg)
	if err != nil {
		return nil, err
	}

	as, err := initAlertSender(ctx, key, publisherSvc)
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

	registryClient, err := ethereum.NewStreamEthClient(ctx, "registry", cfg.Registry.JsonRpc.Url)
	if err != nil {
		return nil, err
	}

	registryService := registry.New(cfg, key.Address, msgClient, registryClient)
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
		health.NewService(ctx, "", healthutils.DefaultHealthServerErrHandler, health.CheckerFrom(
			summarizeReports,
			ethClient, traceClient, blockFeed, txStream, txAnalyzer, blockAnalyzer, agentPool, registryService,
		)),
		txStream,
		txAnalyzer,
		blockAnalyzer,
		scanner.NewScannerAPI(ctx, blockFeed),
		scanner.NewTxLogger(ctx),
		publisherSvc,
	}

	// for performance tests, this flag avoids using registry service
	if !cfg.Registry.Disable {
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
	checkedErr, ok := reports.NameContains("registry.event.checked.error")
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
	blockByNumberTime, ok := reports.NameContains("chain-json-rpc-client.request.block-by-number.time")
	if ok && len(blockByNumberTime.Details) > 0 {
		t, ok := blockByNumberTime.Time()
		if ok {
			checkDelay := time.Since(*t)
			if checkDelay > time.Minute*5 {
				summary.Punc(",")
				summary.Addf("lagging to get block by %d minutes.", int64(checkDelay.Minutes()))
			}
		}
	}
	summary.Punc(".")

	getTxReceiptErr, ok := reports.NameContains("chain-json-rpc-client.request.get-transaction-receipt.error")
	if ok && len(getTxReceiptErr.Details) > 0 {
		summary.Addf("failing to get transaction receipt with error '%s', this can slow down block processing.", getTxReceiptErr.Details)
	}

	traceBlockErr, ok := reports.NameContains("trace-json-rpc-client.request.trace-block.error")
	isTraceBlockNotFoundErr := strings.Contains(traceBlockErr.Details, "not found")
	if ok && len(traceBlockErr.Details) > 0 && !isTraceBlockNotFoundErr {
		summary.Addf("trace api (trace_block) is failing with error '%s'.", traceBlockErr.Details)
		summary.Status(health.StatusFailing)
	}

	return summary.Finish()
}

func Run() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)
	services.ContainerMain("scanner", initServices)
}
