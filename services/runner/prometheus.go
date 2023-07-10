package runner

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/forta-network/forta-core-go/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
func StartPrometheusCollector(serviceHealth ServiceHealth) {
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
		Addr:    ":9107", // TODO: Make configurable
		Handler: mux,
	}

	utils.GoListenAndServe(server)
}

func (pc *promCollector) Describe(ch chan<- *prometheus.Desc) {
}

func (pc *promCollector) Collect(ch chan<- prometheus.Metric) {
	for _, report := range pc.serviceHealth.CheckServiceHealth() {
		if !strings.Contains(report.Name, "service.inspector") {
			continue
		}

		name := strings.Replace(strings.Replace(report.Name, ".", "_", -1), "-", "_", -1)
		val, err := strconv.ParseFloat(report.Details, 64)
		if err != nil {
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("forta", "", name),
				"", nil, nil,
			),
			prometheus.GaugeValue,
			val,
		)
	}
}
