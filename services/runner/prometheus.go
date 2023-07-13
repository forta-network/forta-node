package runner

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

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

// StartPrometheusExporter starts an exporter.
func StartPrometheusCollector(serviceHealth ServiceHealth, port int) {
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

func (pc *promCollector) Collect(ch chan<- prometheus.Metric) {
	for _, report := range pc.serviceHealth.CheckServiceHealth() {
		parts := strings.Split(report.Name, ".service.")
		if len(parts) != 2 {
			continue
		}
		parts = strings.Split(parts[1], ".")
		if len(parts) < 2 {
			continue
		}
		serviceName := toPrometheusName(parts[0])
		reportName := toPrometheusName(strings.Join(parts[1:], "."))

		value := parseReportValue(report)

		metric, err := newPrometheusMetric(value, serviceName, reportName)
		if err != nil {
			logrus.WithError(err).WithField("metric", "forta_"+serviceName+"_"+reportName).Error("failed to create metric")
		}

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

// parseReportValue converts to three types of data to float: timestamp, number, boolean.
// If the value is none of them, finally, it tries to convert the error messages like
// value=1 and label=message.
func parseReportValue(report *health.Report) (value float64) {
	if n, err := strconv.ParseFloat(report.Details, 64); err == nil {
		value = n
		return
	}

	if b, err := strconv.ParseBool(report.Details); err == nil {
		if b {
			value = 1
		}
		return
	}

	if t, err := time.Parse(time.RFC3339, report.Details); err == nil {
		value = float64(t.UTC().Unix())
		return
	}

	// important note: the logic in here is used only if we are trying to use an error message

	switch report.Status {
	case health.StatusOK, health.StatusInfo, health.StatusUnknown:
		value = 0

	case health.StatusFailing, health.StatusLagging, health.StatusDown:
		value = 1
	}
	return
}
