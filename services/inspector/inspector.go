package inspector

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/inspect"
	"github.com/forta-network/forta-core-go/protocol/transform"
	"github.com/forta-network/forta-node/clients"
	"github.com/forta-network/forta-node/clients/messaging"
	"github.com/forta-network/forta-node/config"
	log "github.com/sirupsen/logrus"
)

// Inspector runs continuous inspections.
type Inspector struct {
	ctx context.Context
	cfg InspectorConfig

	msgClient clients.MessageClient

	lastErr health.ErrorTracker

	inspectEvery int
	inspectTrace bool
	inspectCh    chan uint64

	latestBlockEventTime time.Time
}

const blockEventWaitTimeout = time.Minute * 20

type InspectorConfig struct {
	Config    config.Config
	ProxyHost string
	ProxyPort string
}

func (ins *Inspector) Start() error {
	if ins.cfg.Config.LocalModeConfig.Enable && !ins.cfg.Config.LocalModeConfig.EnableInspection {
		log.Warn("inspection is disabled - please enable it from the local mode config if you need it")
		return nil
	}

	ins.registerMessageHandlers()

	go func() {
		for {
			select {
			case <-ins.ctx.Done():
				return

			case <-time.After(blockEventWaitTimeout):
				// if scan api is failing, run a placeholder-like inspection with genesis block
				dialCtx, cancel := context.WithTimeout(ins.ctx, time.Second*3)
				rpcClient, err := rpc.DialContext(dialCtx, ins.cfg.Config.Scan.JsonRpc.Url)
				cancel()
				if err != nil {
					ins.doRunInspection(0)
					continue
				}
				reqCtx, cancel := context.WithTimeout(ins.ctx, time.Second*3)
				blockNum, err := ethclient.NewClient(rpcClient).BlockNumber(reqCtx)
				cancel()
				if err != nil {
					ins.doRunInspection(0)
					continue
				}
				blockNum -= ins.blockNumRemainder(blockNum) // turn it into an expected block num
				ins.doRunInspection(blockNum)

			case blockNum := <-ins.inspectCh:
				ins.doRunInspection(blockNum)
			}
		}
	}()
	return nil
}

func (ins *Inspector) doRunInspection(blockNum uint64) {
	if err := ins.runInspection(ins.ctx, blockNum); err != nil {
		log.WithFields(
			log.Fields{
				"error":             err,
				"inspectingAtBlock": blockNum,
			},
		).Error("finally failed to complete inspection")
	}
}

func (ins *Inspector) runInspection(ctx context.Context, blockNum uint64) error {
	results, err := inspect.Inspect(
		ctx, inspect.InspectionConfig{
			ScanAPIURL:  ins.cfg.Config.Scan.JsonRpc.Url,
			ProxyAPIURL: fmt.Sprintf("http://%s:%s", ins.cfg.ProxyHost, ins.cfg.ProxyPort),
			TraceAPIURL: ins.cfg.Config.Trace.JsonRpc.Url,
			BlockNumber: blockNum,
			CheckTrace:  ins.inspectTrace,
		},
	)
	// publish inspection results even if there are errors, because inspection results are independent from errors
	ins.msgClient.PublishProto(messaging.SubjectInspectionDone, transform.ToProtoInspectionResults(results))

	b, _ := json.Marshal(results)
	log.WithField("results", string(b)).Info("inspection done")
	if err != nil {
		log.WithError(err).Warn("error(s) during inspection")
		return err
	}

	return nil
}

func (ins *Inspector) registerMessageHandlers() {
	ins.msgClient.Subscribe(messaging.SubjectScannerBlock, messaging.ScannerHandler(ins.handleScannerBlock))
}

func (ins *Inspector) handleScannerBlock(payload messaging.ScannerPayload) error {
	ins.latestBlockEventTime = time.Now()

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
	return health.Reports{
		ins.lastErr.GetReport("last-error"),
	}
}

func NewInspector(ctx context.Context, cfg InspectorConfig) (*Inspector, error) {
	msgClient := messaging.NewClient("inspector", fmt.Sprintf("%s:%s", config.DockerNatsContainerName, config.DefaultNatsPort))

	chainSettings := config.GetChainSettings(cfg.Config.ChainID)
	inspectionInterval := chainSettings.InspectionInterval
	if cfg.Config.InspectionConfig.BlockInterval != nil {
		inspectionInterval = *cfg.Config.InspectionConfig.BlockInterval
	}

	inspect.DownloadTestSavingMode = cfg.Config.InspectionConfig.NetworkSavingMode

	return &Inspector{
		ctx:          ctx,
		msgClient:    msgClient,
		cfg:          cfg,
		inspectEvery: inspectionInterval,
		inspectTrace: chainSettings.EnableTrace,
		inspectCh:    make(chan uint64, 1), // let it tolerate being late on one block inspection
	}, nil
}
