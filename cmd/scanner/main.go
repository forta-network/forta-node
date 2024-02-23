package scanner

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	"github.com/forta-network/forta-node/store"

	"github.com/forta-network/forta-core-go/domain"
	"github.com/forta-network/forta-core-go/feeds/timeline"
	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/services/components"
	"github.com/forta-network/forta-node/services/components/estimation"
	"github.com/forta-network/forta-node/services/publisher"
	log "github.com/sirupsen/logrus"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	gethlog "github.com/ethereum/go-ethereum/log"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/ethereum"
	"github.com/forta-network/forta-core-go/feeds"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/scanner"
)

func initTxStream(ctx context.Context, ethClient, traceClient ethereum.Client, cfg config.Config) (*scanner.TxStreamService, feeds.BlockFeed, *estimation.Estimator, estimation.BlockTimeline, error) {
	cfg.Scan.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.JsonRpcProxy.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Scan.JsonRpc.Url)
	cfg.Registry.JsonRpc.Url = utils.ConvertToDockerHostURL(cfg.Registry.JsonRpc.Url)
	cfg.Registry.IPFS.APIURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.APIURL)
	cfg.Registry.IPFS.GatewayURL = utils.ConvertToDockerHostURL(cfg.Registry.IPFS.GatewayURL)

	url := cfg.Scan.JsonRpc.Url
	chainID := config.ParseBigInt(cfg.ChainID)

	if url == "" {
		return nil, nil, nil, nil, fmt.Errorf("scan.jsonRpc.url is required")
	}
	if cfg.Trace.Enabled && cfg.Trace.JsonRpc.Url == "" {
		return nil, nil, nil, nil, fmt.Errorf("trace requires a jsonRpc URL if enabled")
	}

	var rateLimit *time.Ticker
	if cfg.Scan.BlockRateLimit > 0 {
		rateLimit = time.NewTicker(time.Duration(cfg.Scan.BlockRateLimit) * time.Millisecond)
	}

	var maxAgePtr *time.Duration
	// support scanning old block ranges in local mode
	hasLocalModeBlockRange := cfg.LocalModeConfig.Enable && cfg.LocalModeConfig.RuntimeLimits.StopBlock != nil
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
		if runtimeLimits.StartBlock != nil {
			startBlock = big.NewInt(0).SetUint64(*runtimeLimits.StartBlock)
		}
		if runtimeLimits.StopBlock != nil {
			stopBlock = big.NewInt(0).SetUint64(*runtimeLimits.StopBlock)
		}
	}

	if startBlock != nil && stopBlock != nil && !(stopBlock.Cmp(startBlock) > 0) {
		log.Fatal("stop block is not greater than the start block - please check the runtime limits")
	}

	if cfg.Scan.RetryIntervalSeconds > 0 {
		ethClient.SetRetryInterval(time.Second * time.Duration(cfg.Scan.RetryIntervalSeconds))
	} else {
		chainSettings := settings.GetChainSettings(cfg.ChainID)
		ethClient.SetRetryInterval(time.Second * time.Duration(chainSettings.JSONRPCRetryIntervalSeconds))
	}

	blockFeed, err := feeds.NewBlockFeed(ctx, ethClient, traceClient, feeds.BlockFeedConfig{
		ChainID:             chainID,
		Tracing:             cfg.Trace.Enabled,
		RateLimit:           rateLimit,
		SkipBlocksOlderThan: maxAgePtr,
		Start:               startBlock,
		End:                 stopBlock,
	})
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// subscribe to block feed so we can detect block end and trigger exit
	blockErrCh := blockFeed.Subscribe(func(evt *domain.BlockEvent) error {
		return nil
	})

	// subscribe to block feed to construct a timeline and estimate performance
	blockTimeline := timeline.NewBlockTimeline(cfg.ChainID, 5)
	blockFeed.Subscribe(blockTimeline.HandleBlock)
	estimator := estimation.NewEstimator(blockTimeline)

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
		return nil, nil, nil, nil, fmt.Errorf("failed to create the tx stream service: %v", err)
	}

	return txStream, blockFeed, estimator, blockTimeline, nil
}

func initCombinationStream(ctx context.Context, msgClient clients.MessageClient, cfg config.Config) (*scanner.CombinerAlertStreamService, feeds.AlertFeed, error) {
	combinerFeed, err := feeds.NewCombinerFeed(
		ctx, feeds.CombinerFeedConfig{
			APIUrl:            cfg.CombinerConfig.AlertAPIURL,
			Start:             cfg.LocalModeConfig.RuntimeLimits.StartCombiner,
			End:               cfg.LocalModeConfig.RuntimeLimits.StopCombiner,
			CombinerCachePath: cfg.CombinerConfig.CombinerCachePath,
			QueryInterval:     cfg.CombinerConfig.QueryInterval,
		},
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create combiner feed: %v", err)
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

func initTxAnalyzer(
	ctx context.Context, cfg config.Config,
	as clients.AlertSender, stream *scanner.TxStreamService,
	botProcessingComponents components.BotProcessing, msgClient clients.MessageClient,
) (*scanner.TxAnalyzerService, error) {
	return scanner.NewTxAnalyzerService(ctx, scanner.TxAnalyzerServiceConfig{
		TxChannel:     stream.ReadOnlyTxStream(),
		AlertSender:   as,
		MsgClient:     msgClient,
		BotProcessing: botProcessingComponents,
	})
}

func initBlockAnalyzer(
	ctx context.Context, cfg config.Config,
	as clients.AlertSender, stream *scanner.TxStreamService,
	botProcessingComponents components.BotProcessing, msgClient clients.MessageClient,
) (*scanner.BlockAnalyzerService, error) {
	return scanner.NewBlockAnalyzerService(ctx, scanner.BlockAnalyzerServiceConfig{
		BlockChannel:  stream.ReadOnlyBlockStream(),
		AlertSender:   as,
		MsgClient:     msgClient,
		BotProcessing: botProcessingComponents,
	})
}

func initCombinerAlertAnalyzer(
	ctx context.Context, cfg config.Config,
	as clients.AlertSender, stream *scanner.CombinerAlertStreamService,
	botProcessingComponents components.BotProcessing, msgClient clients.MessageClient,
) (*scanner.CombinerAlertAnalyzerService, error) {
	return scanner.NewCombinerAlertAnalyzerService(
		ctx, scanner.CombinerAlertAnalyzerServiceConfig{
			AlertChannel:  stream.ReadOnlyAlertStream(),
			AlertSender:   as,
			MsgClient:     msgClient,
			ChainID:       fmt.Sprintf("%d", cfg.ChainID),
			BotProcessing: botProcessingComponents,
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
	cfg.CombinerConfig.AlertAPIURL = utils.ConvertToDockerHostURL(cfg.CombinerConfig.AlertAPIURL)
	cfg.PublicAPIProxy.Url = utils.ConvertToDockerHostURL(cfg.PublicAPIProxy.Url)
	msgClient := messaging.NewClient("scanner", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	key, err := config.LoadKeyInContainer(cfg)
	if err != nil {
		return nil, err
	}

	ethClient, err := ethereum.NewStreamEthClient(ctx, "chain", cfg.Scan.JsonRpc.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to create stream eth client: %v", err)
	}

	traceClient, err := ethereum.NewStreamEthClient(ctx, "trace", cfg.Trace.JsonRpc.Url)
	if err != nil {
		return nil, fmt.Errorf("failed to create trace stream eth client: %v", err)
	}

	txStream, blockFeed, estimator, blockTimeline, err := initTxStream(ctx, ethClient, traceClient, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create tx stream: %v", err)
	}

	publisherSvc, err := publisher.NewPublisher(ctx, blockTimeline, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %v", err)
	}

	alertSender, err := initAlertSender(ctx, key, publisherSvc, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize alert sender: %v", err)
	}

	var waitBots int
	if cfg.LocalModeConfig.Enable {
		waitBots += len(cfg.LocalModeConfig.BotImages)
		waitBots += len(cfg.LocalModeConfig.Standalone.BotContainers)
		// sharded bots spawn on multiple containers, so total "wait bot" count is shards * target
		for _, bot := range cfg.LocalModeConfig.ShardedBots {
			if bot != nil {
				waitBots += int(bot.Target * bot.Shards)
			}
		}
	}

	botProcessingComponents, err := components.GetBotProcessingComponents(ctx, components.BotProcessingConfig{
		Config:        cfg,
		MessageClient: msgClient,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create bot processing components: %v", err)
	}
	txAnalyzer, err := initTxAnalyzer(ctx, cfg, alertSender, txStream, botProcessingComponents, msgClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize tx analyzer: %v", err)
	}
	blockAnalyzer, err := initBlockAnalyzer(ctx, cfg, alertSender, txStream, botProcessingComponents, msgClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize block analyzer: %v", err)
	}

	// Start the main block feed so all transaction feeds can start consuming.
	if !cfg.Scan.DisableAutostart {
		blockFeed.Start()
	}

	combinationStream, combinationFeed, err := initCombinationStream(ctx, msgClient, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize combiner stream: %v", err)
	}

	combinationAnalyzer, err := initCombinerAlertAnalyzer(ctx, cfg, alertSender, combinationStream, botProcessingComponents, msgClient)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize combiner analyzer: %v", err)
	}

	svcs := []services.Service{
		health.NewService(ctx, "", healthutils.DefaultHealthServerErrHandler, health.CheckerFrom(
			summarizeReports,
			ethClient, traceClient, combinationFeed, blockFeed, txStream,
			txAnalyzer, blockAnalyzer, combinationAnalyzer,
			botProcessingComponents.RequestSender,
			publisherSvc,
			estimator,
		)),
		txStream,
		txAnalyzer,
		blockAnalyzer,
		combinationStream,
		combinationAnalyzer,
		publisherSvc,
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

	summary.Punc(".")

	jsonRpcPerformance, ok := reports.NameContains("json-rpc-performance")
	if ok && jsonRpcPerformance.Status != health.StatusUnknown {
		summary.Addf("scan api performance is estimated as %s (this is different from the SLA score).", jsonRpcPerformance.Details)
	}

	summary.Punc(".")

	jsonRpcDelay, ok := reports.NameContains("json-rpc-delay")
	if ok && jsonRpcPerformance.Status != health.StatusUnknown {
		summary.Addf("the latest block was received %s after creation.", jsonRpcDelay.Details)
	}

	summary.Punc(".")

	return summary.Finish()
}

func isNotFoundErr(errMsg string) bool {
	errMsg = strings.ToLower(errMsg)
	return strings.Contains(errMsg, "not found") || strings.Contains(errMsg, "could not find") ||
		strings.Contains(errMsg, "cannot query unfinalized data") || strings.Contains(errMsg, "unknown block")
}

func Run() {
	gethlog.SetDefault(gethlog.NewLogger(slog.NewTextHandler(os.Stdout, nil)))
	services.ContainerMain("scanner", initServices)
}
