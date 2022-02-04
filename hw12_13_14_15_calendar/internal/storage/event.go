package storage

import "time"

type Event struct {
	ID          string
	OwnerID     string
	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
}
