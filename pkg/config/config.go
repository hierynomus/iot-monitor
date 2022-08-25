package config

import (
	"github.com/hierynomus/iot-monitor/pkg/http"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
)

type Config struct {
	Scraper scraper.Config `yaml:"scraper" viper:"serial"`
	HTTP    http.Config    `yaml:"http" viper:"http"`
}
