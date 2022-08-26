package updater

import (
	"context"
	"sync"

	"github.com/hierynomus/iot-monitor/internal/http"
	"github.com/hierynomus/iot-monitor/internal/prometheus"
	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/hierynomus/iot-monitor/pkg/process"
)

var _ process.Process = (*Updater)(nil)

type Updater struct {
	wg        *sync.WaitGroup
	ch        <-chan iot.RawMessage
	handler   *http.RawMessageHandler
	collector *prometheus.Collector
	converter iot.Converter
}

func NewUpdater(ch <-chan iot.RawMessage, handler *http.RawMessageHandler, collector *prometheus.Collector, converter iot.Converter) *Updater {
	return &Updater{
		wg:        &sync.WaitGroup{},
		ch:        ch,
		handler:   handler,
		converter: converter,
		collector: collector,
	}
}

func (u *Updater) Start(ctx context.Context) error {
	u.wg.Add(1)

	go u.run(ctx)

	return nil
}

func (u *Updater) Stop() error {
	return nil
}

func (u *Updater) Wait() {
	u.wg.Wait()
}

func (u *Updater) run(ctx context.Context) {
	logger := logging.LoggerFor(ctx, "updater")
	defer u.wg.Done()

	for { //nolint:gosimple
		select {
		case t, ok := <-u.ch:
			if !ok {
				logger.Info().Msg("Incoming channel closed, terminating!")
				return
			}

			mm, err := u.converter.Convert(t)
			if err != nil {
				logger.Error().Err(err).Msg("Failed to parse message")
				logger.Trace().Str("raw-message", string(t)).Msg("RawMessage failed to parse")
				continue
			}

			logger.Debug().Msg("Parsed message")

			u.handler.Update(t)
			u.collector.Update(mm)
		}
	}
}
