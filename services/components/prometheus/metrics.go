package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
)

func fqName(metricName string) string {
	return prometheus.BuildFQName("forta", "", metricName)
}

func newGauge(desc *prometheus.Desc, value float64, labels ...string) prometheus.Metric {
	metric, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, value, labels...)
	if err != nil {
		logrus.WithError(err).WithField("metric", desc.String()).Error("failed to create metric")
		return prometheus.NewInvalidMetric(desc, err)
	}
	return metric
}

func transformHealthMetricsToProm(metricKinds []*MetricKind, metrics HealthMetrics) (promMetrics []prometheus.Metric) {
	for _, metricKind := range metricKinds {
		for _, mapping := range metricKind.Mappings {
			if metric, ok := metrics.Get(mapping.FromHealth); ok {
				promMetrics = append(promMetrics, newGauge(metricKind.Desc, metric.Value(), mapping.ToProm))
			}
		}
	}
	return
}

type MetricKind struct {
	Desc     *prometheus.Desc
	Mappings []*Mapping
}

type Mapping struct {
	FromHealth string
	ToProm     string
}

var knownMetricKinds = []*MetricKind{
	////////// input stream

	{
		Desc: prometheus.NewDesc(
			fqName("stream_activity_seconds"), "time elapsed since last input stream activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "tx_stream_event_block_time",
				ToProm:     "tx_stream_block",
			},
			{
				FromHealth: "tx_stream_event_transaction_time",
				ToProm:     "tx_stream_tx",
			},
			// TODO: Figure out later what is going on with these two:
			// {
			// 	FromHealth: "alert_feed_last_alert",
			// 	ToProm: "alert_feed",
			// },
			// {
			// 	FromHealth: "block_feed_last_alert",
			// 	ToProm: "block_feed",
			// },
		},
	},

	////////// JSON-RPC (scan, trace)

	{
		Desc: prometheus.NewDesc(
			fqName("json_rpc_errors"), "detected json-rpc errors",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "chain_json_rpc_client_request_block_by_number_error",
				ToProm:     "block_by_number",
			},
			{
				FromHealth: "trace_json_rpc_client_request_trace_block_error",
				ToProm:     "trace_block",
			},
		},
	},

	{
		Desc: prometheus.NewDesc(
			fqName("json_rpc_activity_seconds"), "time elapsed since last json-rpc activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "chain_json_rpc_client_request_block_by_number_time",
				ToProm:     "block_by_number",
			},
			{
				FromHealth: "trace_json_rpc_client_request_trace_block_time",
				ToProm:     "trace_block",
			},
		},
	},

	////////// updater

	{
		Desc: prometheus.NewDesc(
			fqName("updater_errors"), "detected updater errors",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "updater_event_checked_error",
				ToProm:     "check",
			},
			{
				FromHealth: "updater_event_checked_final_error",
				ToProm:     "final_check",
			},
		},
	},

	{
		Desc: prometheus.NewDesc(
			fqName("updater_activity_seconds"), "time elapsed since last updater activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "updater_event_checked_time",
				ToProm:     "check",
			},
			{
				FromHealth: "updater_event_checked_final_time",
				ToProm:     "final_check",
			},
		},
	},

	////////// analyzer (block, tx, combiner)

	{
		Desc: prometheus.NewDesc(
			fqName("analyzer_activity_seconds"), "time elapsed since last analyzer activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "block_analyzer_event_input_time",
				ToProm:     "block_analyzer_input",
			},
			{
				FromHealth: "block_analyzer_event_output_time",
				ToProm:     "block_analyzer_output",
			},
			{
				FromHealth: "tx_analyzer_event_input_time",
				ToProm:     "tx_analyzer_input",
			},
			{
				FromHealth: "tx_analyzer_event_output_time",
				ToProm:     "tx_analyzer_output",
			},
			{
				FromHealth: "combiner_alert_analyzer_event_input_time",
				ToProm:     "combiner_alert_analyzer_input",
			},
			{
				FromHealth: "combiner_alert_analyzer_event_output_time",
				ToProm:     "combiner_alert_analyzer_output",
			},
		},
	},

	////////// telemetry and bot logs

	{
		Desc: prometheus.NewDesc(
			fqName("telemetry_errors"), "telemetry errors",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "supervisor_event_telemetry_sync_error",
				ToProm:     "telemetry_sync",
			},
			{
				FromHealth: "supervisor_event_custom_telemetry_sync_error",
				ToProm:     "custom_telemetry_sync",
			},
			{
				FromHealth: "supervisor_event_agent_logs_sync_error",
				ToProm:     "bot_logs_sync",
			},
		},
	},

	{
		Desc: prometheus.NewDesc(
			fqName("telemetry_activity_seconds"), "time elapsed since last telemetry activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "supervisor_event_telemetry_sync_time",
				ToProm:     "telemetry_sync",
			},
			{
				FromHealth: "supervisor_event_custom_telemetry_sync_time",
				ToProm:     "custom_telemetry_sync",
			},
			{
				FromHealth: "supervisor_event_agent_logs_sync_time",
				ToProm:     "bot_logs_sync",
			},
		},
	},

	////////// bot registry

	{
		Desc: prometheus.NewDesc(
			fqName("bot_registry_errors"), "detected bot registry errors",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "bot_registry_event_checked_error",
				ToProm:     "check",
			},
		},
	},

	{
		Desc: prometheus.NewDesc(
			fqName("bot_registry_activity_seconds"), "time elapsed since last bot registry activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "bot_registry_event_change_detected_time",
				ToProm:     "change_detected",
			},
			{
				FromHealth: "bot_registry_event_checked_time",
				ToProm:     "check",
			},
		},
	},

	////////// publisher

	{
		Desc: prometheus.NewDesc(
			fqName("publisher_errors"), "detected bot publisher errors",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "publisher_event_batch_publish_error",
				ToProm:     "batch_publish",
			},
		},
	},

	{
		Desc: prometheus.NewDesc(
			fqName("publisher_activity_seconds"), "time elapsed since last publisher activity",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "publisher_event_batch_publish_attempt_time",
				ToProm:     "batch_publish_attempt",
			},
			{
				FromHealth: "publisher_event_batch_skip_time",
				ToProm:     "batch_skip",
			},
			{
				FromHealth: "publisher_event_batch_publish_time",
				ToProm:     "batch_publish",
			},
			{
				FromHealth: "publisher_event_metrics_flush_time",
				ToProm:     "metrics_flush",
			},
		},
	},

	////////// active bots

	{
		Desc: prometheus.NewDesc(
			fqName("active_bots"), "active bot counts",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "sender_agents_total",
				ToProm:     "total",
			},
			{
				FromHealth: "sender_agents_lagging",
				ToProm:     "lagging",
			},
		},
	},

	////////// inspector

	{
		Desc: prometheus.NewDesc(
			fqName("inspector"), "inspection results",
			[]string{"name"}, nil,
		),
		Mappings: inspectorMetricMappings(),
	},

	////////// estimator

	{
		Desc: prometheus.NewDesc(
			fqName("estimator"), "estimation results",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "estimator_json_rpc_performance",
				ToProm:     "json_rpc_performance",
			},
			{
				FromHealth: "estimator_json_rpc_delay",
				ToProm:     "json_rpc_delay",
			},
		},
	},
}

func inspectorMetricMappings() (mappings []*Mapping) {
	for _, healthMetricName := range []string{
		"inspector_api_refs_valid",
		"inspector_expected_score",
		"inspector_last_error",
		"inspector_network_access_outbound",
		"inspector_proxy_api_accessible",
		"inspector_proxy_api_chain_id",
		"inspector_proxy_api_history_support",
		"inspector_proxy_api_is_eth2",
		"inspector_proxy_api_module_eth",
		"inspector_proxy_api_module_net",
		"inspector_proxy_api_module_web3",
		"inspector_registry_api_accessible",
		"inspector_registry_api_assignments",
		"inspector_registry_api_ens",
		"inspector_scan_api_accessible",
		"inspector_scan_api_chain_id",
		"inspector_scan_api_is_eth2",
		"inspector_scan_api_module_eth",
		"inspector_scan_api_module_net",
		"inspector_trace_api_accessible",
		"inspector_trace_api_chain_id",
		"inspector_trace_api_is_eth2",
		"inspector_trace_api_supported",
	} {
		mappings = append(mappings, &Mapping{
			FromHealth: healthMetricName,
			ToProm:     healthMetricName[10:], // drop "inspector_" prefix
		})
	}
	return
}
