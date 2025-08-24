package domain

import (
	"errors"
	"log/slog"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)

func LogErr(err error) slog.Attr {
	return slog.Any("error", err)
}
