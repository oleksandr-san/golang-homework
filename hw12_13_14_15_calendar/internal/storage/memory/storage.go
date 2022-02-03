package memorystorage

import (
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

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// check if event time is free?
	// get events by day, select events with this hour,
	// how to do that with SQL?
	//   a) index event by begin and end timestamp!
	//   b) search events, that during created event. If found -> error that this time is busy

	s.events[event.ID] = event
	return nil
}

func (s *Storage) GetEvent(id string) (*storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event, ok := s.events[id]; ok {
		return &event, nil
	}

	return nil, storage.ErrNotFound
}

func (s *Storage) UpdateEvent(id string, updatedEvent storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if storedEvent, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	} else if storedEvent.ID != updatedEvent.ID {
		return storage.ErrInvalidEvent
	}

	s.events[id] = updatedEvent
	return nil
}

func (s *Storage) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrNotFound
	}

	delete(s.events, id)
	return nil
}

func (s *Storage) ListEventsForDay(day time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]
		if event.Start.Year() == day.Year() && event.Start.Month() == day.Month() && event.Start.Day() == day.Day() {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) ListEventsForWeek(firstDay time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]

		eventYear, eventWeek := event.Start.ISOWeek()
		year, week := firstDay.ISOWeek()
		if eventYear == year && eventWeek == week {
			events = append(events, event)
		}
	}

	return events, nil
}

func (s *Storage) ListEventsForMonth(firstDay time.Time) ([]storage.Event, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	events := []storage.Event{}

	for id := range s.events {
		event := s.events[id]

		if event.Start.Year() != firstDay.Year() ||
			event.Start.Month() != firstDay.Month() {
			continue
		}

		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Start.Before(events[j].Start)
	})

	return events, nil
}
