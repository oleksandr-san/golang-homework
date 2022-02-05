package memorystorage

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]storage.Event
}

func New() *Storage {
	return &Storage{
		mu:     sync.RWMutex{},
		events: map[string]storage.Event{},
	}
}

func (s *Storage) CreateEvent(_ context.Context, event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event
	return nil
}

func (s *Storage) ReadEvent(_ context.Context, eventID, _ string) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, ok := s.events[eventID]; ok {
		return &event, nil
	}

	return nil, storage.ErrNotFound
}

func (s *Storage) UpdateEvent(_ context.Context, updatedEvent storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[updatedEvent.ID]; !ok {
		return storage.ErrNotFound
	}

	s.events[updatedEvent.ID] = updatedEvent
	return nil
}

func (s *Storage) DeleteEvent(_ context.Context, eventID, _ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[eventID]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, eventID)
	return nil
}

func (s *Storage) ListEventsForDay(_ context.Context, _ string, day time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]
		if event.StartTime.Year() == day.Year() &&
			event.StartTime.Month() == day.Month() &&
			event.StartTime.Day() == day.Day() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) ListEventsForWeek(_ context.Context, _ string, firstDay time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]

		eventYear, eventWeek := event.StartTime.ISOWeek()
		year, week := firstDay.ISOWeek()
		if eventYear == year && eventWeek == week {
			events = append(events, event)
		}
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events, nil
}

func (s *Storage) ListEventsForMonth(_ context.Context, _ string, firstDay time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]

		if event.StartTime.Year() != firstDay.Year() ||
			event.StartTime.Month() != firstDay.Month() {
			continue
		}

		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].StartTime.Before(events[j].StartTime)
	})

	return events, nil
}
