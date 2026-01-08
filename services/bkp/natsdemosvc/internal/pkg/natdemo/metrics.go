package natdemo

import (
	"github.com/prometheus/client_golang/prometheus"
)

type NATDemoMetrics struct {
	RequestsTotal *prometheus.CounterVec
}

func NewNATDemoMetrics(namespace string, subsystem string) *NATDemoMetrics {
	return &NATDemoMetrics{
		RequestsTotal: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "requests_total",
				Help:      "Total number of requests processed by natdemo",
			},
			[]string{"operation"},
		),
	}
}

// Collect implements the metric collection callback.
func (m *NATDemoMetrics) Collect(ch chan<- prometheus.Metric) {
	m.RequestsTotal.Collect(ch)
}
