package storage

import "errors"

var (
	ErrDateBusy     = errors.New("another event is assigned to this date")
	ErrNotFound     = errors.New("event not found")
	ErrInvalidEvent = errors.New("event invalid")
)
