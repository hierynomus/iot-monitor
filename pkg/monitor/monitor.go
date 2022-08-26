package monitor

import (
	"context"
	"os"
	"sync"
	"syscall"

	"github.com/hierynomus/iot-monitor/internal/http"
	"github.com/hierynomus/iot-monitor/internal/prometheus"
	"github.com/hierynomus/iot-monitor/internal/updater"
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/iot-monitor/pkg/converter"
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
	"github.com/ztrue/shutdown"
)

type Monitor struct {
	config    config.MonitorConfig
	scraper   scraper.Scraper
	server    *http.Server
	updater   *updater.Updater
	WaitGroup *sync.WaitGroup
}

func NewMonitor(ctx context.Context, config config.MonitorConfig, scraper scraper.Scraper, provider exporter.MetricProvider, converter converter.Converter) (*Monitor, error) {
	collector := prometheus.NewCollector(provider)

	server := http.NewServer(ctx, config.HTTP())

	handler := http.NewRawMessageHandler(ctx, config.RawMessageContentType())
	server.AddHandler("/", handler)
	promhandler, err := prometheus.NewPrometheusHandler(ctx, collector)
	if err != nil {
		return nil, err
	}

	server.AddHandler("/metrics", promhandler)

	updater := updater.NewUpdater(scraper.Output(), handler, collector, converter)

	return &Monitor{
		config:    config,
		scraper:   scraper,
		updater:   updater,
		server:    server,
		WaitGroup: &sync.WaitGroup{},
	}, nil
}

func (m *Monitor) Start(ctx context.Context) error {
	shutdown.AddWithParam(func(s os.Signal) {
		logger := logging.LoggerFor(ctx, "shutdown-hook")
		logger.Warn().Str("signal", s.String()).Msg("Received signal, shutting down")
		if err := m.scraper.Stop(); err != nil {
			logger.Error().Err(err).Msg("Error stopping scraper, abort graceful shutdown")
			return
		}

		m.scraper.Wait()
		m.updater.Wait()
		if err := m.server.Stop(); err != nil {
			logger.Error().Err(err).Msg("Failed to gracefully stop server")
		}

		m.server.Wait()

		logger.Info().Msg("All processes stopped, terminating!")
	})

	if err := m.scraper.Start(ctx); err != nil {
		return err
	}

	if err := m.updater.Start(ctx); err != nil {
		return err
	}

	if err := m.server.Start(ctx); err != nil {
		return err
	}

	shutdown.Listen(syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	return nil
}
