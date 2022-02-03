package storage

import "time"

type Event struct {
	ID           string
	Title        string
	Description  string
	Start        time.Time
	Duration     time.Duration
	OwnerID      string
	Notification time.Duration
}
