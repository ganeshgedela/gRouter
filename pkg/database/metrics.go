package database

import (
	"database/sql"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// MetricsCollector holds the Prometheus metrics for database stats
type MetricsCollector struct {
	dbName string
	db     *sql.DB

	openConnections  *prometheus.GaugeVec
	idleConnections  *prometheus.GaugeVec
	inUseConnections *prometheus.GaugeVec
	waitCount        *prometheus.GaugeVec
	waitDuration     *prometheus.GaugeVec
}

// NewMetricsCollector creates a new collector for the given database
func NewMetricsCollector(dbName string, db *sql.DB) *MetricsCollector {
	m := &MetricsCollector{
		dbName: dbName,
		db:     db,
		openConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "db_open_connections",
			Help: "The number of established connections both in use and idle.",
		}, []string{"db_name"}),
		idleConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "db_idle_connections",
			Help: "The number of idle connections.",
		}, []string{"db_name"}),
		inUseConnections: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "db_in_use_connections",
			Help: "The number of connections currently in use.",
		}, []string{"db_name"}),
		waitCount: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "db_wait_count",
			Help: "The total number of connections waited for.",
		}, []string{"db_name"}),
		waitDuration: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "db_wait_duration_seconds",
			Help: "The total time blocked waiting for a new connection.",
		}, []string{"db_name"}),
	}

	// Register metrics with the global registry
	prometheus.MustRegister(m.openConnections)
	prometheus.MustRegister(m.idleConnections)
	prometheus.MustRegister(m.inUseConnections)
	prometheus.MustRegister(m.waitCount)
	prometheus.MustRegister(m.waitDuration)

	return m
}

// Start begins collecting metrics in the background
func (m *MetricsCollector) Start(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			stats := m.db.Stats()

			m.openConnections.WithLabelValues(m.dbName).Set(float64(stats.OpenConnections))
			m.idleConnections.WithLabelValues(m.dbName).Set(float64(stats.Idle))
			m.inUseConnections.WithLabelValues(m.dbName).Set(float64(stats.InUse))
			m.waitCount.WithLabelValues(m.dbName).Set(float64(stats.WaitCount))
			m.waitDuration.WithLabelValues(m.dbName).Set(stats.WaitDuration.Seconds())
		}
	}()
}
