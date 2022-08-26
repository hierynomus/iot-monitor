package http

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/hierynomus/iot-monitor/pkg/iot"
	"github.com/hierynomus/iot-monitor/pkg/logging"
)

var _ http.Handler = (*RawMessageHandler)(nil)

type RawMessageHandler struct {
	ctx            context.Context
	mutex          sync.RWMutex
	LastUpdateTime time.Time
	RawMessage     iot.RawMessage
	ContentType    string
}

func NewRawMessageHandler(ctx context.Context, contentType string) *RawMessageHandler {
	return &RawMessageHandler{
		ctx:         ctx,
		ContentType: contentType,
	}
}

func (h *RawMessageHandler) Update(msg iot.RawMessage) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.RawMessage = msg
	h.LastUpdateTime = time.Now()
}

func (h *RawMessageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	logger := logging.LoggerFor(h.ctx, "raw-message-handler")
	w.Header().Add("Content-Type", fmt.Sprintf("%s; charset=utf-8", h.ContentType))
	w.Header().Add("Last-Modified", h.LastUpdateTime.Format(http.TimeFormat))
	if _, err := w.Write([]byte(h.RawMessage)); err != nil {
		logger.Error().Err(err).Msg("Failed to write raw message")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
