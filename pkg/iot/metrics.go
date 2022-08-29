package iot

import (
	iotprom "github.com/hierynomus/iot-monitor/pkg/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

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

func NewCounterL(namespace, name, help string, labels map[string]string) prometheus.Counter {
	return iotprom.NewCounter(prometheus.CounterOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        help,
		ConstLabels: prometheus.Labels(labels),
	})
}

func NewGaugeL(namespace, name, help string, labels map[string]string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   namespace,
		Name:        name,
		Help:        help,
		ConstLabels: prometheus.Labels(labels),
	})
}

func Labels(kvs ...string) map[string]string {
	m := make(map[string]string)

	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return m
}
