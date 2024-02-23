package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/inspect"
	"github.com/forta-network/forta-core-go/inspect/scorecalc"
	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

const (
	blockEventWaitTimeout            = time.Minute * 20
	defaultPublishInspectionInterval = time.Second * 5
)

// Inspector runs continuous inspections.
type Inspector struct {
	ctx context.Context
	cfg InspectorConfig

	msgClient clients.MessageClient

	lastErr          health.ErrorTracker
	indicatorReports []health.Report
	trackerMu        sync.RWMutex

	latestInspection          *inspect.InspectionResults
	latestInspectionMu        sync.RWMutex
	inspectionPublishInterval time.Duration

	inspectEvery int
	inspectTrace bool
	inspectCh    chan uint64
}

type InspectorConfig struct {
	Config         config.Config
	ProxyHost      string
	ProxyPort      string
	ScannerAddress string
}

func (ins *Inspector) Start() error {
	if ins.cfg.Config.LocalModeConfig.Enable && !ins.cfg.Config.LocalModeConfig.ForceEnableInspection {
		log.Warn("inspection is disabled - please enable it from the local mode config using 'forceEnableInspection' if you need it")
		return nil
	}

	if *ins.cfg.Config.InspectionConfig.InspectAtStartup {
		blockNumber := ins.getClosestBlockToInspect()
		ins.runInspection(blockNumber)
	}

	ins.registerMessageHandlers()

	go ins.inspectionPublisher(ins.ctx)

	go func() {
		for {
			select {
			case <-ins.ctx.Done():
				return

			case <-time.After(blockEventWaitTimeout):
				blockNum := ins.getClosestBlockToInspect()
				ins.runInspection(blockNum)

			case blockNum := <-ins.inspectCh:
				ins.runInspection(blockNum)
			}
		}
	}()
	return nil
}

func (ins *Inspector) inspectionPublisher(ctx context.Context) error {
	t := time.NewTicker(ins.inspectionPublishInterval)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-t.C:
			ins.latestInspectionMu.RLock()
			if ins.latestInspection != nil {
				ins.msgClient.PublishProto(messaging.SubjectInspectionDone, inspect.ToProtoInspectionResults(ins.latestInspection))
			}
			ins.latestInspectionMu.RUnlock()
		}
	}
}

func (ins *Inspector) runInspection(blockNum uint64) error {
	inspectCtx, cancel := context.WithTimeout(ins.ctx, time.Minute*2)
	results, err := inspect.Inspect(
		inspectCtx, inspect.InspectionConfig{
			ScanAPIURL:         ins.cfg.Config.Scan.JsonRpc.Url,
			ProxyAPIURL:        fmt.Sprintf("http://%s:%s", ins.cfg.ProxyHost, ins.cfg.ProxyPort),
			TraceAPIURL:        ins.cfg.Config.Trace.JsonRpc.Url,
			BlockNumber:        blockNum,
			CheckTrace:         ins.inspectTrace,
			RegistryAPIURL:     ins.cfg.Config.Registry.JsonRpc.Url,
			ENSContractAddress: ins.cfg.Config.ENSConfig.ContractAddress,
			ScannerAddress:     ins.cfg.ScannerAddress,
		},
	)

	if ins.cfg.Config.JsonRpcProxy.JsonRpc.Url != "" {
		results.Inputs.ProxyAPIURL = ins.cfg.Config.JsonRpcProxy.JsonRpc.Url
	} else {
		results.Inputs.ProxyAPIURL = ins.cfg.Config.Scan.JsonRpc.Url
	}

	cancel()

	// use inspection results even if there are errors
	// because inspection results are independent of errors
	ins.latestInspectionMu.Lock()
	ins.latestInspection = results
	ins.latestInspectionMu.Unlock()

	ins.trackerMu.Lock()
	ins.indicatorReports = nil
	for indicatorName, indicatorValue := range results.Indicators {
		ins.indicatorReports = append(
			ins.indicatorReports, health.Report{
				Name:    indicatorName,
				Status:  health.StatusInfo,
				Details: strconv.FormatFloat(indicatorValue, 'f', -1, 64),
			},
		)
	}
	chainID := uint64(ins.cfg.Config.ChainID)
	inspectionScore, _ := scorecalc.NewScoreCalculator([]scorecalc.ScoreCalculatorConfig{{ChainID: chainID}}).CalculateScore(chainID, results)
	ins.indicatorReports = append(
		ins.indicatorReports, health.Report{
			Name:    "expected-score",
			Status:  health.StatusInfo,
			Details: strconv.FormatFloat(inspectionScore, 'f', -1, 64),
		},
	)
	ins.trackerMu.Unlock()

	b, _ := json.Marshal(results)
	log.WithFields(
		log.Fields{
			"results":           string(b),
			"inspectingAtBlock": blockNum,
		},
	).Info("inspection done")

	if err != nil {
		log.WithFields(
			log.Fields{
				"error":             err,
				"inspectingAtBlock": blockNum,
			},
		).Error("failed to execute inspection")

		return err
	}

	return nil
}

func (ins *Inspector) registerMessageHandlers() {
	ins.msgClient.Subscribe(messaging.SubjectScannerBlock, messaging.ScannerHandler(ins.handleScannerBlock))
}

func (ins *Inspector) handleScannerBlock(payload messaging.ScannerPayload) error {
	if payload.LatestBlockInput > 0 && ins.blockNumRemainder(payload.LatestBlockInput) == 0 {
		// inspect from N blocks back to avoid synchronizations issues
		inspectionBlockNum := payload.LatestBlockInput - uint64(ins.inspectEvery)
		logger := log.WithFields(
			log.Fields{
				"triggeredAtBlock":  payload.LatestBlockInput,
				"inspectingAtBlock": inspectionBlockNum,
			},
		)
		logger.Info("triggering inspection")
		// non-blocking insert
		select {
		case ins.inspectCh <- inspectionBlockNum:
			logger.Info("successfully triggered new inspection")
		default:
			logger.Info("failed to trigger new inspection: already busy")
		}
	}
	return nil
}

func (ins *Inspector) blockNumRemainder(blockNum uint64) uint64 {
	return blockNum % uint64(ins.inspectEvery)
}

func (ins *Inspector) Stop() error {
	return nil
}

func (ins *Inspector) Name() string {
	return "inspector"
}

// Health implements health.Reporter interface.
func (ins *Inspector) Health() health.Reports {
	reports := health.Reports{
		ins.lastErr.GetReport("last-error"),
	}
	ins.trackerMu.RLock()
	for _, report := range ins.indicatorReports {
		reportCopy := report
		reports = append(reports, &reportCopy)
	}
	ins.trackerMu.RUnlock()

	return reports
}

func (ins *Inspector) getClosestBlockToInspect() uint64 {
	// if scan api is failing, run a placeholder-like inspection with genesis block
	dialCtx, cancel := context.WithTimeout(ins.ctx, time.Second*3)
	rpcClient, err := rpc.DialContext(dialCtx, ins.cfg.Config.Scan.JsonRpc.Url)
	cancel()
	if err != nil {
		return 0
	}

	reqCtx, cancel := context.WithTimeout(ins.ctx, time.Second*3)
	blockNum, err := ethclient.NewClient(rpcClient).BlockNumber(reqCtx)
	cancel()
	if err != nil {
		return 0
	}

	blockNum -= ins.blockNumRemainder(blockNum) // turn it into an expected block num

	return blockNum
}

func NewInspector(ctx context.Context, cfg InspectorConfig) (*Inspector, error) {
	msgClient := messaging.NewClient("inspector", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	chainSettings := settings.GetChainSettings(cfg.Config.ChainID)
	inspectionInterval := chainSettings.InspectionInterval
	if cfg.Config.InspectionConfig.BlockInterval != nil {
		inspectionInterval = *cfg.Config.InspectionConfig.BlockInterval
	}

	publishInterval := defaultPublishInspectionInterval
	if cfg.Config.Publish.Batch.IntervalSeconds != nil {
		publishInterval = (time.Duration(*cfg.Config.Publish.Batch.IntervalSeconds) * time.Second) / 3
	}

	return &Inspector{
		ctx:                       ctx,
		msgClient:                 msgClient,
		cfg:                       cfg,
		inspectEvery:              inspectionInterval,
		inspectTrace:              chainSettings.EnableTrace,
		inspectCh:                 make(chan uint64, 1), // let it tolerate being late on one block inspection
		inspectionPublishInterval: publishInterval,
	}, nil
}
