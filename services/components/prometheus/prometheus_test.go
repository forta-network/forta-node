package prometheus

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
	}
}

var testMetricKinds = []*MetricKind{
	{
		Desc: prometheus.NewDesc(
			fqName("service_bar"), "",
			[]string{"name"}, nil,
		),
		Mappings: []*Mapping{
			{
				FromHealth: "bar_m1",
				ToProm:     "m1",
			},
		},
	},
}

func TestPrometheusCollector(t *testing.T) {
	r := require.New(t)

	collector := &promCollector{
		serviceHealth: &testServiceHealth{},
		metricKinds:   testMetricKinds,
	}

	descCh := make(chan *prometheus.Desc, 1)
	collector.Describe(descCh)
	close(descCh)
	r.Len(descCh, 0)

	metricCh := make(chan prometheus.Metric, 1)
	collector.Collect(metricCh)
	close(metricCh)

	var allMetrics []prometheus.Metric
	for metric := range metricCh {
		allMetrics = append(allMetrics, metric)
	}

	r.Len(allMetrics, 1)

	var encodedMetric io_prometheus_client.Metric

	m1 := allMetrics[0]
	r.NoError(m1.Write(&encodedMetric))
	r.Equal(float64(1420687622144), *encodedMetric.Gauge.Value)
	label := encodedMetric.GetLabel()[0]
	r.Equal("name", label.GetName())
	r.Equal("m1", label.GetValue())
}

func TestStartPrometheusCollector(t *testing.T) {
	StartCollector(&testServiceHealth{}, testMetricKinds, 9107)
}
