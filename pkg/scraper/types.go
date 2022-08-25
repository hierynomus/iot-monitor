package scraper

import (
	"context"

	"github.com/hierynomus/iot-monitor/pkg/process"
)

type RawMessage string

type Config interface {
	RawMessageContentType() string
	Get(key string) interface{}
	GetString(key string) string
	GetInt(key string) int
	GetBool(key string) bool
	Validate() error
}

var _ process.Process = (Scraper)(nil) // compile-time interface check

type Scraper interface {
	Start(ctx context.Context) error
	Stop() error
	Wait() error
	Output() <-chan RawMessage
}
