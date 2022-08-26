package iotmonitor

import (
	"context"

	"github.com/hierynomus/iot-monitor/cmd"
	"github.com/hierynomus/iot-monitor/pkg/monitor"
	"github.com/rs/zerolog"
)

func StartMonitor(ctx context.Context, config interface{}, monitorStarter func() (*monitor.Monitor, error)) error {
	c := cmd.RootCommand(config)
	c.AddCommand(cmd.StartCommand(monitorStarter))

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return c.ExecuteContext(ctx)
}
