package httpDelivery

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Lovodia/Product/internal/domain"
)

func renderError(w http.ResponseWriter, r *http.Request, err error) {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrBadRequest),
		errors.Is(err, domain.ErrInvalidImput):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
	default:
		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(domain.ErrorResponse{
		Error: err.Error(),
	})

	slog.Error("http error",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote", r.RemoteAddr),
		slog.Int("status", status),
		slog.Any("error", err),
	)
}
