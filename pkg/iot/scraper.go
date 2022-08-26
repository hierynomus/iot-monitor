package iot

import (
	"context"

	"github.com/hierynomus/iot-monitor/pkg/process"
)

type RawMessage string

var _ process.Process = (Scraper)(nil) // compile-time interface check

type Scraper interface {
	Start(ctx context.Context) error
	Stop() error
	Wait()
	Output() <-chan RawMessage
}
