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
		if len(parts) != 2 {
			continue
		}
		serviceName := parts[0]
		reportName := strings.Join(parts[1:], ".")

		metricName := reportName
		metricName = strings.ReplaceAll(metricName, "-", "_")
		metricName = strings.ReplaceAll(metricName, ".", "_")
		metricName = fmt.Sprintf("forta_%s_%s", serviceName, metricName)

		// converting three types of data to float: timestamp, number, boolean
		// finally, converting the error messages like value=1 and label=message
		var (
			value  float64
			labels []string
		)
		if t, err := time.Parse(time.RFC3339, report.Details); err != nil {
			value = float64(t.UTC().Unix())
		} else if n, err := strconv.ParseFloat(report.Details, 64); err != nil {
			value = n
		} else if b, err := strconv.ParseBool(report.Details); err != nil {
			if b {
				value = 1
			}
		} else {

			// important note: the logic in here is used only if we are trying to use an error message

			if report.Status != health.StatusOK {
				labels = append(labels, string(report.Status), report.Details)
			}

			switch report.Status {
			case health.StatusOK, health.StatusInfo, health.StatusUnknown:
				value = 0

			case health.StatusFailing, health.StatusLagging:
				value = 1

			case health.StatusDown:
				value = -1
			}

		}

		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName("forta", serviceName, metricName),
				"", nil, nil,
			),
			prometheus.GaugeValue,
			value,
		)
	}
}
