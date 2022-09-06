package prometheus

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestShouldSetCounter(t *testing.T) {
	counter := NewCounter(prometheus.CounterOpts{
		Name: "test",
		Help: "test",
	})
	counter.Set(42)
	assert.Equal(t, 42.0, *counter.(*smartCounter).value)
	dto := &dto.Metric{}
	assert.NoError(t, counter.Write(dto))
	assert.Equal(t, 42.0, dto.Counter.GetValue())
}

func TestShouldSetCounterTwiceIncrementsUnderlyingCounter(t *testing.T) {
	counter := NewCounter(prometheus.CounterOpts{
		Name: "test",
		Help: "test",
	})
	counter.Set(42)
	counter.Set(43)
	assert.Equal(t, 43.0, *counter.(*smartCounter).value)
	dto := &dto.Metric{}
	assert.NoError(t, counter.Write(dto))
	assert.Equal(t, 43.0, dto.Counter.GetValue())
}

func TestShouldNotBeAbleToDecrementCounterUsingSet(t *testing.T) {
	counter := NewCounter(prometheus.CounterOpts{
		Name: "test",
		Help: "test",
	})
	counter.Set(42)
	counter.Set(41)
	assert.Equal(t, 42.0, *counter.(*smartCounter).value)
	dto := &dto.Metric{}
	assert.NoError(t, counter.Write(dto))
	assert.Equal(t, 42.0, dto.Counter.GetValue())
}
