package sqlstorage

import (
	"context"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	driverName     string
	dataSourceName string
	db             *sqlx.DB
}

func New(driverName, dataSourceName string) *Storage {
	return &Storage{
		driverName:     driverName,
		dataSourceName: dataSourceName,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect(s.driverName, s.dataSourceName)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(storage.Event) error {
	return nil
}

func (s *Storage) UpdateEvent(id string, event storage.Event) error {
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	return nil
}

func (s *Storage) ListEventsForDay(day time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (s *Storage) ListEventsForWeek(firstDay time.Time) ([]storage.Event, error) {
	return nil, nil
}

func (s *Storage) ListEventsForMonth(firstDay time.Time) ([]storage.Event, error) {
	return nil, nil
}
