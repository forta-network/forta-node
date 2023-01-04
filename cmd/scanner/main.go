package scanner

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/forta-network/forta-node/store"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/services/publisher"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-core-go/security"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/registry"
	"github.com/forta-network/forta-node/services/scanner"
	"github.com/forta-network/forta-node/services/scanner/agentpool"
)

func initTxStream(ctx context.Context, ethClient, traceClient ethereum.Client, cfg config.Config) (*scanner.TxStreamService, feeds.BlockFeed, error) {
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.JsonRpcProxy.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
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

	var maxAgePtr *time.Duration
	// support scanning old block ranges in local mode
	hasLocalModeBlockRange := cfg.LocalModeConfig.Enable && cfg.LocalModeConfig.RuntimeLimits.StopBlock > 0
	if !hasLocalModeBlockRange && cfg.Scan.BlockMaxAgeSeconds > 0 {
		maxAge := time.Duration(cfg.Scan.BlockMaxAgeSeconds) * time.Second
		maxAgePtr = &maxAge
	}

	var (
		startBlock *big.Int
		stopBlock  *big.Int
	)
	if cfg.LocalModeConfig.Enable {
		runtimeLimits := cfg.LocalModeConfig.RuntimeLimits
		if runtimeLimits.StartBlock > 0 {
			startBlock = big.NewInt(0).SetUint64(runtimeLimits.StartBlock)
		}
		if runtimeLimits.StopBlock > 0 {
			stopBlock = big.NewInt(0).SetUint64(runtimeLimits.StopBlock)
		}
	}

	ethClient.SetRetryInterval(time.Second * time.Duration(cfg.Scan.RetryIntervalSeconds))

	blockFeed, err := feeds.NewBlockFeed(ctx, ethClient, traceClient, feeds.BlockFeedConfig{
		ChainID:             chainID,
		Tracing:             cfg.Trace.Enabled,
		RateLimit:           rateLimit,
		SkipBlocksOlderThan: maxAgePtr,
		Offset:              getBlockOffset(cfg),
		Start:               startBlock,
		End:                 stopBlock,
	})
	if err != nil {
		return nil, nil, err
	}

	// subscribe to block feed so we can detect block end and trigger exit
	blockErrCh := blockFeed.Subscribe(func(evt *domain.BlockEvent) error {
		return nil
	})
	// detect end block, wait for scanning to finish, trigger exit
	go func() {
		err := <-blockErrCh
		if err == nil {
			return
		}

		if err == context.Canceled {
			return
		}

		if err != feeds.ErrEndBlockReached {
			log.WithError(err).Panic("unexpected failure in block feed")
		}
		log.Info("end block reached - triggering exit")
		var delay time.Duration
		if cfg.LocalModeConfig.Enable {
			delay = time.Duration(cfg.LocalModeConfig.RuntimeLimits.StopTimeoutSeconds) * time.Second
		}
		services.TriggerExit(delay)
	}()

	txStream, err := scanner.NewTxStreamService(ctx, ethClient, blockFeed, scanner.TxStreamServiceConfig{
		JsonRpcConfig:       cfg.Scan.JsonRpc,
		TraceJsonRpcConfig:  cfg.Trace.JsonRpc,
		SkipBlocksOlderThan: maxAgePtr,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the tx stream service: %v", err)
	}

	return txStream, blockFeed, nil
}

// getBlockOffset either returns the default offset configured for the chain or
// the safe offset if required.
func getBlockOffset(cfg config.Config) int {
	chainSettings := settings.GetChainSettings(cfg.ChainID)

	if cfg.AdvancedConfig.SafeOffset {
		return chainSettings.SafeOffset
	}

	scanURL := strings.Trim(cfg.Scan.JsonRpc.Url, "/")
	proxyURL := strings.Trim(cfg.JsonRpcProxy.JsonRpc.Url, "/")
	if len(proxyURL) > 0 && proxyURL != scanURL {
		return chainSettings.SafeOffset
	}

	return chainSettings.DefaultOffset
}

func initCombinationStream(ctx context.Context, msgClient *messaging.Client, cfg config.Config) (*scanner.CombinerAlertStreamService, feeds.AlertFeed, error) {
	combinerFeed, err := feeds.NewCombinerFeed(
		ctx, feeds.CombinerFeedConfig{
			APIUrl:            cfg.CombinerConfig.AlertAPIURL,
			Start:             cfg.LocalModeConfig.RuntimeLimits.StartCombiner,
			End:               cfg.LocalModeConfig.RuntimeLimits.StopCombiner,
			CombinerCachePath: cfg.CombinerConfig.CombinerCachePath,
		},
	)
	if err != nil {
		return nil, nil, err
	}

	combinerStream, err := scanner.NewCombinerAlertStreamService(
		ctx, combinerFeed, msgClient, scanner.CombinerAlertStreamServiceConfig{
			Start: cfg.LocalModeConfig.RuntimeLimits.StartCombiner,
			End:   cfg.LocalModeConfig.RuntimeLimits.StopCombiner,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create the tx stream service: %v", err)
	}

	// subscribe to combiner feed so we can detect combiner stop and trigger exit
	combinerErrCh := combinerFeed.RegisterHandler(
		func(evt *domain.AlertEvent) error {
			return nil
		},
	)

	// detect end time, wait for scanning to finish, trigger exit
	go func() {
		err := <-combinerErrCh
		if err == nil {
			return
		}

		if err == context.Canceled {
			return
		}

		if err != feeds.ErrCombinerStopReached {
			log.WithError(err).Panic("unexpected failure in block feed")
		}
		log.Info("combiner stop reached - triggering exit")

		var delay time.Duration
		if cfg.LocalModeConfig.Enable {
			delay = time.Duration(cfg.LocalModeConfig.RuntimeLimits.StopTimeoutSeconds) * time.Second
		}
		services.TriggerExit(delay)
	}()

	return combinerStream, combinerFeed, nil
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

func initCombinerAlertAnalyzer(ctx context.Context, cfg config.Config, as clients.AlertSender, stream *scanner.CombinerAlertStreamService, ap *agentpool.AgentPool, msgClient clients.MessageClient) (*scanner.CombinerAlertAnalyzerService, error) {
	return scanner.NewCombinerAlertAnalyzerService(
		ctx, scanner.CombinerAlertAnalyzerServiceConfig{
			AlertChannel: stream.ReadOnlyAlertStream(),
			AlertSender:  as,
			AgentPool:    ap,
			MsgClient:    msgClient,
			ChainID:      fmt.Sprintf("%d", cfg.ChainID),
		},
	)
}

func initAlertSender(ctx context.Context, key *keystore.Key, pubClient clients.PublishClient, cfg config.Config) (clients.AlertSender, error) {
	ds, err := store.NewDeduplicationStore(cfg)
	if err != nil {
		return nil, err
	}
	return clients.NewAlertSender(ctx, pubClient, clients.AlertSenderConfig{
		Key: key,
		DS:  ds,
	})
}

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	// can't dial localhost - need to dial host gateway from container
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.Trace.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Trace.JsonRpc.Url)
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)
	cfg.Registry.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.APIURL)
	cfg.Registry.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.GatewayURL)
	cfg.Publish.APIURL = utils.ConvertToDockerHostURL(cfg.Publish.APIURL)
	cfg.Publish.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Publish.IPFS.APIURL)
	cfg.Publish.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Publish.IPFS.GatewayURL)
	cfg.LocalModeConfig.WebhookURL = utils.ConvertToDockerHostURL(cfg.LocalModeConfig.WebhookURL)
	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := security.LoadKey(config.DefaultContainerKeyDirPath)
	if err != nil {
		return nil, err
	}

	publisherSvc, err := publisher.NewPublisher(ctx, cfg)
	if err != nil {
		return nil, err
	}

	as, err := initAlertSender(ctx, key, publisherSvc, cfg)
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

	combinationStream, combinationFeed, err := initCombinationStream(ctx, msgClient, cfg)
	if err != nil {
		return nil, err
	}

	registryClient, err := ethereum.NewStreamEthClient(ctx, "registry", cfg.Registry.JsonRpc.Url)
	if err != nil {
		return nil, err
	}
	registryService := registry.New(cfg, key.Address, msgClient, registryClient)

	var waitBots int
	if cfg.LocalModeConfig.Enable {
		waitBots = len(cfg.LocalModeConfig.BotImages)
	}

	agentPool := agentpool.NewAgentPool(ctx, cfg.Scan, msgClient, waitBots)
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, as, txStream, agentPool, msgClient)
	if err != nil {
		return nil, err
	}

	combinationAnalyzer, err := initCombinerAlertAnalyzer(ctx, cfg, as, combinationStream, agentPool, msgClient)
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
			ethClient, traceClient, combinationFeed, blockFeed, txStream, txAnalyzer, blockAnalyzer, combinationAnalyzer, agentPool, registryService,
			publisherSvc,
		)),
		txStream,
		txAnalyzer,
		blockAnalyzer,
		combinationStream,
		combinationAnalyzer,
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

	// report block request failures but ignore "not found"s because we hit them when we are
	// asking for the latest block that is not just yet available
	blockByNumberErr, ok := reports.NameContains("chain-json-rpc-client.request.block-by-number.error")
	if ok && len(blockByNumberErr.Details) > 0 && !isNotFoundErr(blockByNumberErr.Details) {
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
	if ok && len(traceBlockErr.Details) > 0 && !isNotFoundErr(traceBlockErr.Details) {
		summary.Addf("trace api (trace_block) is failing with error '%s'.", traceBlockErr.Details)
		summary.Status(health.StatusFailing)
	}
	summary.Punc(".")

	batchPublishErr, ok := reports.NameContains("publisher.event.batch-publish.error")
	if ok && len(batchPublishErr.Details) > 0 {
		summary.Addf("failed to publish the last batch with error '%s'", batchPublishErr.Details)
		summary.Status(health.StatusFailing)
	}

	return summary.Finish()
}

func isNotFoundErr(errMsg string) bool {
	return strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "could not find") ||
		strings.Contains(errMsg, "cannot query unfinalized data")
}

func Run() {
	gethlog.Root().SetHandler(gethlog.StdoutHandler)
	services.ContainerMain("scanner", initServices)
}
