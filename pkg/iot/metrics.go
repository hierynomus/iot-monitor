package iot

import "github.com/prometheus/client_golang/prometheus"

type MetricMessage map[string]Metric

type Metric struct {
	Value string
	Unit  string
}

type MetricCollector interface {
	prometheus.Collector
	prometheus.Metric
}

type MetricProvider interface {
	Metrics() map[string]MetricCollector
}
