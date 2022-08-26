package cmd

import (
	"github.com/hierynomus/iot-monitor/pkg/monitor"
	"github.com/hierynomus/iot-monitor/version"
	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

func StartCommand(monitorStarter func() (*monitor.Monitor, error)) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start the SMA Monitor",
		Long:  "Start the SMA Monitor",
		RunE:  RunStart(monitorStarter),
	}
}

func RunStart(monitorStarter func() (*monitor.Monitor, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		logger := log.Ctx(cmd.Context())
		logger.Info().Str("version", version.Version).Str("commit", version.Commit).Str("date", version.Date).Msg("Starting P1 Monitor")

		m, err := monitorStarter()
		if err != nil {
			return err
		}

		return m.Start(cmd.Context())
	}
}
