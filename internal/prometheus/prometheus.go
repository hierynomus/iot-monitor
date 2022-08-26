package prometheus

import (
	"context"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewPrometheusHandler(ctx context.Context, collector *Collector) (http.Handler, error) {
	reg := prometheus.NewRegistry()

	if err := collector.RegisterMetrics(context.Background(), reg); err != nil {
		return nil, err
	}

	return promhttp.HandlerFor(reg, promhttp.HandlerOpts{}), nil
}
