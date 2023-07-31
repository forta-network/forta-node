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
}

type promHTTPLogger struct{}

func (l promHTTPLogger) Println(v ...interface{}) {
	logrus.Error(v...)
}

// StartCollector starts a collector.
func StartCollector(serviceHealth ServiceHealth, port int) {
	prometheus.MustRegister(version.NewCollector("forta_node"))

	var collector prometheus.Collector = &promCollector{serviceHealth: serviceHealth}
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

type healthMetric struct {
	MetricName string
	Report     *health.Report
}

func (metric *healthMetric) Value() float64 {
	return parseReportValue(metric.Report)
}

func (metric *healthMetric) SinceUnixTs() float64 {
	return sinceUnixTs(metric.Report.Details)
}

type healthMetrics []*healthMetric

func (metrics healthMetrics) Get(name string) (*healthMetric, bool) {
	for _, metric := range metrics {
		if metric.MetricName == name {
			return metric, true
		}
	}
	return nil, false
}

func (pc *promCollector) Collect(ch chan<- prometheus.Metric) {
	var healthMetrics healthMetrics
	for _, report := range pc.serviceHealth.CheckServiceHealth() {
		parts := strings.Split(report.Name, ".service.")
		if len(parts) != 2 {
			continue
		}
		subParts := strings.Split(parts[1], ".")
		if len(subParts) < 2 {
			continue
		}

		healthMetrics = append(healthMetrics, &healthMetric{
			MetricName: toPrometheusName(parts[1]),
			Report:     report,
		})
	}

	sendMetrics(ch, transformHealthMetricsToProm(healthMetrics)...)
}

func sendMetrics(ch chan<- prometheus.Metric, metrics ...prometheus.Metric) {
	for _, metric := range metrics {
		ch <- metric
	}
}

func newPrometheusMetric(value float64, serviceName, reportName string) (prometheus.Metric, error) {
	desc := prometheus.NewDesc(
		prometheus.BuildFQName("forta", serviceName, reportName),
		"", nil, nil,
	)
	metric, err := prometheus.NewConstMetric(desc, prometheus.GaugeValue, value)
	if err != nil {
		metric = prometheus.NewInvalidMetric(desc, err)
	}
	return metric, err
}

func toPrometheusName(name string) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	return name
}
