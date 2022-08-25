package cmd

import (
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/monitor"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
	"github.com/hierynomus/iot-monitor/version"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

func StartCommand(cfg *config.Config, scraper scraper.Scraper, provider exporter.MetricProvider) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the SMA Monitor",
		Long:  "Start the SMA Monitor",
		RunE:  RunStart(cfg, scraper, provider),
	}
}

func RunStart(cfg *config.Config, scraper scraper.Scraper, provider exporter.MetricProvider) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger := log.Ctx(cmd.Context())
		logger.Info().Str("version", version.Version).Str("commit", version.Commit).Str("date", version.Date).Msg("Starting P1 Monitor")

		m, err := monitor.NewMonitor(cmd.Context(), cfg, scraper, provider)
		if err != nil {
			return err
		}

		return m.Start(cmd.Context())
	}
}
