package domain

import (
	"errors"
)

var (
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)
