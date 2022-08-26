package prometheus

import (
	"context"
	"strconv"
	"sync"

	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/prometheus/client_golang/prometheus"
)

var _ prometheus.Collector = (*Collector)(nil)

type Collector struct {
	lock    sync.RWMutex
	metrics map[string]iot.MetricCollector
}

// NewCollector returns a new Collector.
func NewCollector(provider iot.MetricProvider) *Collector {
	return &Collector{
		metrics: provider.Metrics(),
	}
}

func (c *Collector) RegisterMetrics(ctx context.Context, reg prometheus.Registerer) error {
	for _, m := range c.metrics {
		if err := reg.Register(m); err != nil {
			return err
		}
	}
	return nil
}

func (c *Collector) Update(msg iot.MetricMessage) {
	c.lock.Lock()
	defer c.lock.Unlock()
	logger := logging.LoggerFor(context.Background(), "prometheus-collector")

	for k, metric := range c.metrics {
		if v, ok := msg[k]; ok {
			fl, err := strconv.ParseFloat(v.Value, 64) //nolint:gomnd
			if err != nil {
				logger.Warn().Err(err).Str("key", k).Str("value", v.Value).Msg("Failed to parse value to float")
				continue
			}
			switch m := metric.(type) {
			case prometheus.Gauge:
				m.Set(fl)
			case prometheus.Counter:
				m.Add(fl)
			default:
				logger.Warn().Err(err).Str("key", k).Str("value", v.Value).Msg("Unknown prometheus metric type")
			}
		}
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, metric := range c.metrics {
		ch <- metric.Desc()
	}
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	for _, metric := range c.metrics {
		ch <- metric
	}
}
