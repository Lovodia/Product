package domain

import "errors"

var (
	ErrBadRequest   = errors.New("bad request")
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidImput = errors.New("invalid imput")
	ErrConflict     = errors.New("conflict")
	ErrInternal     = errors.New("internal error")
)
