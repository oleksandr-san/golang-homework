package app

import (
	"context"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	storage Storage
	logger  Logger
}

type Logger interface {
	Info(msg string)
	Debug(msg string)
	Warning(msg string)
	Error(msg string)
}

type Storage interface { // TODO
	CreateEvent(storage.Event) error
	UpdateEvent(id string, event storage.Event) error
	DeleteEvent(id string) error
	ListEventsForDay(day time.Time) ([]storage.Event, error)
	ListEventsForWeek(firstDay time.Time) ([]storage.Event, error)
	ListEventsForMonth(firstDay time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{storage: storage, logger: logger}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
