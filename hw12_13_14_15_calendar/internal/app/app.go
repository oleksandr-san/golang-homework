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

type Storage interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	ReadEvent(ctx context.Context, eventID, ownerID string) (*storage.Event, error)
	UpdateEvent(ctx context.Context, event storage.Event) error
	DeleteEvent(ctx context.Context, eventID, ownerID string) error
	ListEventsForDay(ctx context.Context, ownerID string, day time.Time) ([]storage.Event, error)
	ListEventsForWeek(ctx context.Context, ownerID string, firstDay time.Time) ([]storage.Event, error)
	ListEventsForMonth(ctx context.Context, ownerID string, firstDay time.Time) ([]storage.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{storage: storage, logger: logger}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	return a.storage.CreateEvent(ctx, storage.Event{ID: id, Title: title})
}
