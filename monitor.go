package iotmonitor

import (
	"context"

	"github.com/hierynomus/iot-monitor/cmd"
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
	"github.com/rs/zerolog"
)

func StartMonitor(ctx context.Context, scraper scraper.Scraper, provider exporter.MetricProvider) error {
	cfg := &config.Config{}
	c := cmd.RootCommand(cfg)
	c.AddCommand(cmd.StartCommand(cfg, scraper, provider))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return c.ExecuteContext(ctx)
}
