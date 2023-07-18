package runner

import (
	"testing"

	"github.com/forta-network/forta-core-go/clients/health"
	"github.com/prometheus/client_golang/prometheus"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/require"
)

type testServiceHealth struct{}

func (*testServiceHealth) CheckServiceHealth() (allReports health.Reports) {
	return health.Reports{
		{
			Name:    "foo.service.bar.m1",
			Status:  health.StatusOK,
			Details: "1420687622144",
		},
		{
			Name:    "foo.service.bar.m2",
			Status:  health.StatusOK,
			Details: "true",
		},
		{
			Name:    "foo.service.bar.m3",
			Status:  health.StatusOK,
			Details: "2020-10-11T01:02:03Z",
		},
		{
			Name:    "foo.service.bar.m4",
			Status:  health.StatusOK,
			Details: "message 1",
		},
		{
			Name:    "foo.service.bar.m5",
			Status:  health.StatusFailing,
			Details: "message 2",
		},
		{
			Name:    "skipped-invalid-name-pattern",
			Status:  health.StatusFailing,
			Details: "message X",
		},
		{
			Name:    "foo.service.skipped-single-name",
			Status:  health.StatusFailing,
			Details: "message Y",
		},
	}
}

func TestPrometheusCollector(t *testing.T) {
	r := require.New(t)

	collector := &promCollector{
		serviceHealth: &testServiceHealth{},
	}

	descCh := make(chan *prometheus.Desc, 1)
	collector.Describe(descCh)
	close(descCh)
	r.Len(descCh, 0)

	metricCh := make(chan prometheus.Metric, 5)
	collector.Collect(metricCh)
	close(metricCh)

	var allMetrics []prometheus.Metric
	for metric := range metricCh {
		allMetrics = append(allMetrics, metric)
	}

	var encodedMetric io_prometheus_client.Metric

	m1 := allMetrics[0]
	r.NoError(m1.Write(&encodedMetric))
	r.Equal(float64(1420687622144), *encodedMetric.Gauge.Value)

	m2 := allMetrics[1]
	r.NoError(m2.Write(&encodedMetric))
	r.Equal(float64(1), *encodedMetric.Gauge.Value)

	m3 := allMetrics[2]
	r.NoError(m3.Write(&encodedMetric))
	r.Equal(float64(1602378123), *encodedMetric.Gauge.Value)

	m4 := allMetrics[3]
	r.NoError(m4.Write(&encodedMetric))
	r.Equal(float64(0), *encodedMetric.Gauge.Value)

	m5 := allMetrics[4]
	r.NoError(m5.Write(&encodedMetric))
	r.Equal(float64(1), *encodedMetric.Gauge.Value)
}

func TestStartPrometheusCollector(t *testing.T) {
	StartPrometheusCollector(&testServiceHealth{}, 9107)
}
