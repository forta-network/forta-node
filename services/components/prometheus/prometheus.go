package prometheus

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

// ServiceHealth checks the service health.
type ServiceHealth interface {
	CheckServiceHealth() (allReports health.Reports)
}

type promCollector struct {
	serviceHealth ServiceHealth
	metricKinds   []*MetricKind
}

type promHTTPLogger struct{}

func (l promHTTPLogger) Println(v ...interface{}) {
	logrus.Error(v...)
}

// StartCollector starts a collector.
func StartCollector(serviceHealth ServiceHealth, metricKinds []*MetricKind, port int) {
	prometheus.MustRegister(version.NewCollector("forta_node"))

	if metricKinds == nil {
		metricKinds = knownMetricKinds
	}

	var collector prometheus.Collector = &promCollector{
		serviceHealth: serviceHealth,
		metricKinds:   metricKinds,
	}
	prometheus.MustRegister(collector)

	mux := http.NewServeMux()

	mux.HandleFunc("/-/healthy", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})
	mux.HandleFunc("/-/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "OK")
	})

	mux.Handle("/metrics", promhttp.InstrumentMetricHandler(
		prometheus.DefaultRegisterer,
		promhttp.HandlerFor(
			prometheus.DefaultGatherer,
			promhttp.HandlerOpts{
				ErrorLog: &promHTTPLogger{},
			},
		),
	))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}
	logrus.WithField("port", port).Info("starting prometheus server")

	utils.GoListenAndServe(server)
}

func (pc *promCollector) Describe(ch chan<- *prometheus.Desc) {
}

type HealthMetric struct {
	MetricName string
	Report     *health.Report
}

func (metric *HealthMetric) Value() float64 {
	return parseReportValue(metric.Report)
}

type HealthMetrics []*HealthMetric

func (metrics HealthMetrics) Get(name string) (*HealthMetric, bool) {
	for _, metric := range metrics {
		if metric.MetricName == name {
			return metric, true
		}
	}
	return nil, false
}

func (pc *promCollector) Collect(ch chan<- prometheus.Metric) {
	var healthMetrics HealthMetrics
	for _, report := range pc.serviceHealth.CheckServiceHealth() {
		parts := strings.Split(report.Name, ".service.")
		if len(parts) != 2 {
			continue
		}
		subParts := strings.Split(parts[1], ".")
		if len(subParts) < 2 {
			continue
		}

		healthMetrics = append(healthMetrics, &HealthMetric{
			MetricName: toPrometheusName(parts[1]),
			Report:     report,
		})
	}

	sendMetrics(ch, transformHealthMetricsToProm(pc.metricKinds, healthMetrics)...)
}

func sendMetrics(ch chan<- prometheus.Metric, metrics ...prometheus.Metric) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func toPrometheusName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}
