package inspector

import (
	"context"
	"strconv"
	"strings"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/inspect"
	"github.com/forta-network/forta-core-go/inspect/scorecalc"
	"github.com/forta-network/forta-core-go/protocol/settings"
	"github.com/forta-network/forta-node/config"
	"github.com/forta-network/forta-node/healthutils"
	"github.com/forta-network/forta-node/services"
	"github.com/forta-network/forta-node/services/inspector"
)

var nodeConfig config.Config

func initServices(ctx context.Context, cfg config.Config) ([]services.Service, error) {
	nodeConfig = cfg

	inspector, err := inspector.NewInspector(ctx, inspector.InspectorConfig{
		Config:    cfg,
		ProxyHost: config.DockerJSONRPCProxyContainerName,
		ProxyPort: config.DefaultJSONRPCProxyPort,
	})
	if err != nil {
		return nil, err
	}

	return []services.Service{
		health.NewService(
			ctx, "", healthutils.DefaultHealthServerErrHandler,
			health.CheckerFrom(summarizeReports, inspector),
		),
		inspector,
	}, nil
}

func summarizeReports(reports health.Reports) *health.Report {
	summary := health.NewSummary()

	chainSetings := settings.GetChainSettings(nodeConfig.ChainID)

	var failingApis []string

	scanAccessible, ok := reports.NameContains(inspect.IndicatorScanAPIAccessible)
	if !ok {
		return summary.Finish()
	}
	if ok && scanAccessible.Details != "1" {
		failingApis = append(failingApis, "scan")
	}
	traceAccessible, ok := reports.NameContains(inspect.IndicatorTraceAccessible)
	if ok && traceAccessible.Details != "1" && chainSetings.EnableTrace {
		failingApis = append(failingApis, "trace")
	}

	if len(failingApis) > 0 {
		summary.Addf("something is wrong with %s api.", strings.Join(failingApis, ", "))
		summary.Status(health.StatusFailing)
	}

	var incompatibleApis []string

	expectedChainID := strconv.FormatInt(int64(nodeConfig.ChainID), 10)

	scanChainID, ok := reports.NameContains(inspect.IndicatorScanAPIChainID)
	if ok && scanChainID.Details != expectedChainID {
		incompatibleApis = append(incompatibleApis, "scan")
	}
	traceChainID, ok := reports.NameContains(inspect.IndicatorTraceAPIChainID)
	if ok && traceChainID.Details != expectedChainID && chainSetings.EnableTrace {
		incompatibleApis = append(incompatibleApis, "trace")
	}
	if len(incompatibleApis) > 0 {
		summary.Addf("different chain detected from %s api (expected chain id %s).", strings.Join(incompatibleApis, ", "), expectedChainID)
		summary.Status(health.StatusFailing)
	}

	traceSupported, ok := reports.NameContains(inspect.IndicatorTraceSupported)
	if ok && traceSupported.Details == "-1" && chainSetings.EnableTrace {
		summary.Add("trace api does not support `trace_block`.")
		summary.Status(health.StatusFailing)
	}

	totalMemory, ok := reports.NameContains(inspect.IndicatorResourcesMemoryTotal)
	if ok {
		mem, _ := strconv.ParseFloat(totalMemory.Details, 64)
		if mem < scorecalc.MinTotalMemoryRequired {
			summary.Add("low total memory.")
			summary.Status(health.StatusFailing)
		}
	}

	return summary.Finish()
}

func Run() {
	services.ContainerMain("inspector", initServices)
}
