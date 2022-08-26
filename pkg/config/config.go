package config

import (
	"time"
)

type HTTPConfig struct {
	ListenAddress string        `yaml:"listen-address" viper:"listen-address" env:"LISTEN_ADDRESS" default:":8080"`
	Timeout       time.Duration `yaml:"timeout" viper:"timeout" env:"TIMEOUT" default:"30s"`
}

type MonitorConfig interface {
	HTTP() HTTPConfig
	RawMessageContentType() string
}
