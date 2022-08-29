package prometheus

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

var _ prometheus.Collector = (*smartCounter)(nil)
var _ prometheus.Metric = (*smartCounter)(nil)
var _ SettableCounter = (*smartCounter)(nil)

type SettableCounter interface {
	prometheus.Counter
	Set(float64)
}

func NewCounter(opts prometheus.CounterOpts) SettableCounter {
	return &smartCounter{
		counter: prometheus.NewCounter(opts),
		value:   nil,
		mutex:   sync.RWMutex{},
	}
}

type smartCounter struct {
	counter prometheus.Counter
	value   *float64
	mutex   sync.RWMutex
}

func (c *smartCounter) Inc() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.counter.Inc()
}

func (c *smartCounter) Add(delta float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.counter.Add(delta)
}

func (c *smartCounter) Set(value float64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.value == nil {
		c.counter.Add(value)
		v := &value
		c.value = v
	} else {
		delta := value - *c.value
		if delta > 0 {
			c.counter.Add(delta)
			c.value = &value
		}

	}

}

func (c *smartCounter) Collect(ch chan<- prometheus.Metric) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.counter.Collect(ch)
}

func (c *smartCounter) Describe(ch chan<- *prometheus.Desc) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	c.counter.Describe(ch)
}

func (c *smartCounter) Desc() *prometheus.Desc {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.counter.Desc()
}

func (c *smartCounter) Write(dto *dto.Metric) error {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.counter.Write(dto)
}
