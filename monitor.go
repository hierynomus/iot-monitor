package iotmonitor

import (
	"context"
	"fmt"

	"github.com/hierynomus/autobind"
	"github.com/hierynomus/iot-monitor/cmd"
	"github.com/hierynomus/iot-monitor/pkg/monitor"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

type Bootstrapper struct {
	Name        string
	Description string
	EnvPrefix   string
	Config      interface{}
	Starter     func() (*monitor.Monitor, error)
	Binder      *autobind.Autobinder
	Viper       *viper.Viper
}

func NewBootstrapper(name, description, envPrefix string, config interface{}, starter func() (*monitor.Monitor, error)) *Bootstrapper {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.AddConfigPath(".")
	vp.AddConfigPath(fmt.Sprintf("/etc/%s", name))
	vp.SetConfigType("yaml")

	return &Bootstrapper{
		Name:        name,
		Description: description,
		Config:      config,
		Starter:     starter,
		Binder:      &autobind.Autobinder{UseNesting: true, EnvPrefix: envPrefix, ConfigObject: config, Viper: vp, SetDefaults: true},
		Viper:       vp,
	}
}

func (b *Bootstrapper) Start(ctx context.Context) error {
	c := cmd.RootCommand(b.Config, b.Name, b.Description, b.Binder, b.Viper)

	c.AddCommand(cmd.StartCommand(b.Starter))
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return c.ExecuteContext(ctx)
}
