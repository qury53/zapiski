package core_http_response

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	core_errors "github.com/qury53/GC_learn/api/internal/core/errors"
	"github.com/qury53/GC_learn/platform/logger"
)

type HTTPResponseHandler struct {
	log *logger.Logger
	rw  http.ResponseWriter
}

func NewHTTPResponseHandler(log *logger.Logger, rw http.ResponseWriter) *HTTPResponseHandler {
	return &HTTPResponseHandler{log: log, rw: rw}
}

func (h *HTTPResponseHandler) JSONResponse(responseBody any, statusCode int) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)
	if err := json.NewEncoder(h.rw).Encode(&responseBody); err != nil && h.log != nil {
		h.log.Error("write HTTP response", slog.Any("error", err))
	}
}

func (h *HTTPResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)
	if h.log != nil {
		h.log.Error(msg, slog.Any("error", err))
	}
	h.errorResponse(statusCode, err, msg)
}

func (h *HTTPResponseHandler) ErrorResponse(err error, msg string) {
	var statusCode int
	level := slog.LevelError

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		level = slog.LevelWarn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		level = slog.LevelDebug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		level = slog.LevelWarn
	default:
		statusCode = http.StatusInternalServerError
	}

	if h.log != nil {
		attrs := []any{slog.Any("error", err)}
		h.log.Log(contextBackground(), level, msg, attrs...)
	}
	h.errorResponse(statusCode, err, msg)
}

func contextBackground() context.Context { return context.Background() }

func (h *HTTPResponseHandler) errorResponse(statusCode int, err error, msg string) {
	h.rw.Header().Set("Content-Type", "application/json")
	h.rw.WriteHeader(statusCode)
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	_ = json.NewEncoder(h.rw).Encode(map[string]string{"message": msg, "error": errMsg})
}
