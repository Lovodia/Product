package httpdelivery

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/Lovodia/Product/internal/domain"
)

func renderJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		slog.Error("failed to encode JSON", domain.LogErr(err))
	}
}

func renderError(w http.ResponseWriter, r *http.Request, err error) error {
	status := http.StatusInternalServerError

	switch {
	case errors.Is(err, domain.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, domain.ErrBadRequest),
		errors.Is(err, domain.ErrInvalidInput):
		status = http.StatusBadRequest
	case errors.Is(err, domain.ErrConflict):
		status = http.StatusConflict
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if encodeErr := json.NewEncoder(w).Encode(domain.ErrorResponse{
		Error: err.Error(),
	}); encodeErr != nil {
		slog.Error("failed to encode error JSON response", domain.LogErr(encodeErr))
	}

	slog.Error("http error",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
		slog.String("remote", r.RemoteAddr),
		slog.Int("status", status),
		domain.LogErr(err),
	)

	return err
}
