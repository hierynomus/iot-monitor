package cmd

import (
	"fmt"
	"os"

	"github.com/hierynomus/autobind"
	"github.com/hierynomus/iot-monitor/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	VerboseFlag      = "verbose"
	VerboseFlagShort = "v"
)

func RootCommand(cfg *config.Config) *cobra.Command {
	var verbosity int

	cmd := &cobra.Command{
		Use:   "iot-monitor",
		Short: "P1 Monitor",
		Long:  "P1 Monitor",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			switch verbosity {
			case 0:
				// Nothing to do
			case 1:
				zerolog.SetGlobalLevel(zerolog.DebugLevel)
			default:
				zerolog.SetGlobalLevel(zerolog.TraceLevel)
			}

			vp := viper.New()
			vp.SetConfigName("config")
			vp.AddConfigPath(".")
			vp.AddConfigPath("/etc/iot-monitor")
			vp.SetConfigType("yaml")

			logger := log.Ctx(cmd.Context())

			if err := vp.ReadInConfig(); err != nil {
				if _, ok := err.(viper.ConfigFileNotFoundError); ok {
					logger.Warn().Msg("No config file found... Continuing with defaults")
					// Config file not found; ignore error if desired
				} else {
					fmt.Printf("%s", err)
					os.Exit(1)
				}
			}

			binder := &autobind.Autobinder{
				UseNesting:   true,
				EnvPrefix:    "K2P",
				ConfigObject: cfg,
				Viper:        vp,
				SetDefaults:  true,
			}

			binder.Bind(cmd.Context(), cmd, []string{})

			return nil
		},
	}

	cmd.PersistentFlags().CountVarP(&verbosity, VerboseFlag, VerboseFlagShort, "Print verbose logging to the terminal (use multiple times to increase verbosity)")

	return cmd
}
