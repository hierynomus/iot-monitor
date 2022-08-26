package converter

import (
	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
)

type Converter interface {
	Convert(in scraper.RawMessage) (exporter.MetricMessage, error)
}
