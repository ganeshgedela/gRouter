package nats

import (
	"context"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Optimised buckets for NATS (sub-millisecond to seconds)
	// 100µs, 500µs, 1ms, 5ms, 10ms, 25ms, 50ms, 100ms, 250ms, 500ms, 1s, 2.5s, 5s, 10s
	natsBuckets = []float64{0.0001, 0.0005, 0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

	// Metrics for publishers
	publishCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messaging_publish_total",
		Help: "Total number of messages published",
	}, []string{"subject", "type", "status"})

	publishDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "messaging_publish_duration_seconds",
		Help:    "Duration of message publishing in seconds",
		Buckets: natsBuckets,
	}, []string{"subject", "type"})

	// Metrics for requests (separate from fire-and-forget publishes)
	requestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messaging_request_total",
		Help: "Total number of requests sent",
	}, []string{"subject", "type", "status"})

	requestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "messaging_request_duration_seconds",
		Help:    "Duration of request-reply latency in seconds",
		Buckets: natsBuckets,
	}, []string{"subject", "type"})

	// Metrics for subscribers
	subscribeCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "messaging_subscribe_total",
		Help: "Total number of messages received",
	}, []string{"subject", "type", "status", "service"})

	subscribeDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "messaging_subscribe_duration_seconds",
		Help:    "Duration of message processing in seconds",
		Buckets: natsBuckets,
	}, []string{"subject", "type", "service"})
)

// --- Metrics Middleware ---

// MetricsMiddleware returns a middleware that tracks message processing metrics
func MetricsMiddleware() SubscriberMiddleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(ctx context.Context, subject string, env *MessageEnvelope) error {
			start := time.Now()
			err := next(ctx, subject, env)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			service := env.Source
			if service == "" {
				service = "unknown"
			}

			subscribeCounter.WithLabelValues(subject, env.Type, status, service).Inc()
			subscribeDuration.WithLabelValues(subject, env.Type, service).Observe(duration.Seconds())

			return err
		}
	}
}

// PublisherMetricsMiddleware returns a middleware that tracks message publishing metrics
func PublisherMetricsMiddleware() PublisherMiddleware {
	return func(next PublisherFunc) PublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) error {
			start := time.Now()
			err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			publishCounter.WithLabelValues(subject, msgType, status).Inc()
			publishDuration.WithLabelValues(subject, msgType).Observe(duration.Seconds())

			return err
		}
	}
}

// JSPublisherMetricsMiddleware returns a middleware that tracks JetStream publishing metrics
func JSPublisherMetricsMiddleware() JSPublisherMiddleware {
	return func(next JSPublisherFunc) JSPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (*nats.PubAck, error) {
			start := time.Now()
			ack, err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			publishCounter.WithLabelValues(subject, msgType, status).Inc()
			publishDuration.WithLabelValues(subject, msgType).Observe(duration.Seconds())

			return ack, err
		}
	}
}

// JSAsyncPublisherMetricsMiddleware returns a middleware that tracks async JetStream publishing metrics
func JSAsyncPublisherMetricsMiddleware() JSAsyncPublisherMiddleware {
	return func(next JSAsyncPublisherFunc) JSAsyncPublisherFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, opts *PublishOptions) (nats.PubAckFuture, error) {
			start := time.Now()
			future, err := next(ctx, subject, msgType, data, opts)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			publishCounter.WithLabelValues(subject, msgType, status).Inc()
			publishDuration.WithLabelValues(subject, msgType).Observe(duration.Seconds())

			return future, err
		}
	}
}

// RequestMetricsMiddleware returns a middleware that tracks request metrics
func RequestMetricsMiddleware() RequestMiddleware {
	return func(next RequestFunc) RequestFunc {
		return func(ctx context.Context, subject string, msgType string, data interface{}, timeout time.Duration, opts *PublishOptions) (*MessageEnvelope, error) {
			start := time.Now()
			resp, err := next(ctx, subject, msgType, data, timeout, opts)
			duration := time.Since(start)

			status := "success"
			if err != nil {
				status = "error"
			}

			// Use specific request metrics
			requestCounter.WithLabelValues(subject, msgType, status).Inc()
			requestDuration.WithLabelValues(subject, msgType).Observe(duration.Seconds())

			return resp, err
		}
	}
}
