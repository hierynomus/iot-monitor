package exporter

import "github.com/prometheus/client_golang/prometheus"

type MetricCollector interface {
	prometheus.Collector
	prometheus.Metric
}

type MetricProvider interface {
	Metrics() map[string]MetricCollector
}
