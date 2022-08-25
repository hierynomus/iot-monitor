package monitor

import (
	"context"
	"sync"

	"github.com/hierynomus/iot-monitor/pkg/exporter"
	"github.com/hierynomus/iot-monitor/pkg/http"
	"github.com/hierynomus/iot-monitor/pkg/logging"
	"github.com/hierynomus/iot-monitor/pkg/process"
	"github.com/hierynomus/iot-monitor/pkg/scraper"
	"github.com/roaldnefs/go-dsmr"
)

var _ process.Process = (*Updater)(nil)

type Updater struct {
	WaitGroup *sync.WaitGroup
	ch        <-chan scraper.RawMessage
	handler   *http.RawMessageHandler
	collector *exporter.Collector
}

func NewUpdater(ch <-chan scraper.RawMessage, handler *http.RawMessageHandler, collector *exporter.Collector) *Updater {
	return &Updater{
		WaitGroup: &sync.WaitGroup{},
		ch:        ch,
		handler:   handler,
		collector: collector,
	}
}

func (u *Updater) Start(ctx context.Context) error {
	u.WaitGroup.Add(1)

	go u.run(ctx)

	return nil
}

func (u *Updater) Stop() error {
	return nil
}

func (u *Updater) Wait() error {
	u.WaitGroup.Wait()
	return nil
}

func (u *Updater) run(ctx context.Context) {
	logger := logging.LoggerFor(ctx, "updater")
	defer u.WaitGroup.Done()

	for { //nolint:gosimple
		select {
		case t, ok := <-u.ch:
			if !ok {
				logger.Info().Msg("Updater channel closed, terminating!")
				return
			}

			parsedTelegram, err := dsmr.ParseTelegram(string(t))
			if err != nil {
				logger.Error().Err(err).Msg("Failed to parse telegram")
				continue
			}

			logger.Debug().Msg("Parsed telegram")

			u.handler.Update(t)
			u.collector.Update(parsedTelegram)
		}
	}
}
