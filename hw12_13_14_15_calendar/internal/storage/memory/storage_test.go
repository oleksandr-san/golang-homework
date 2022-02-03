package memorystorage

import (
	"testing"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("describe existing event succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "Test"}

		err := s.CreateEvent(testEvent)
		require.NoError(t, err)

		event, err := s.GetEvent("1")
		require.NoError(t, err)

		require.Equal(t, testEvent, *event)
	})

	t.Run("describe unexisting event fails", func(t *testing.T) {
		s := New()
		event, err := s.GetEvent("1")

		require.Nil(t, event)
		require.Equal(t, err, storage.ErrNotFound)
	})

	t.Run("update existing even succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "A"}
		err := s.CreateEvent(testEvent)
		require.NoError(t, err)

		updatedEvent := storage.Event{ID: "1", Title: "B"}
		err = s.UpdateEvent("1", updatedEvent)
		require.NoError(t, err)

		storedEvent, err := s.GetEvent("1")
		require.NoError(t, err)
		require.NotNil(t, storedEvent)
		require.Equal(t, updatedEvent, *storedEvent)
	})

	t.Run("update unexising event fails", func(t *testing.T) {
		s := New()
		err := s.UpdateEvent("1", storage.Event{})
		require.Equal(t, err, storage.ErrNotFound)
	})

	t.Run("update existing event ID fails", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "A"}
		err := s.CreateEvent(testEvent)
		require.NoError(t, err)

		updatedEvent := storage.Event{ID: "2", Title: "B"}
		err = s.UpdateEvent("1", updatedEvent)
		require.Equal(t, err, storage.ErrInvalidEvent)
	})

	t.Run("delete existing event succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "A"}
		err := s.CreateEvent(testEvent)
		require.NoError(t, err)

		err = s.DeleteEvent(testEvent.ID)
		require.NoError(t, err)

		storedEvent, err := s.GetEvent(testEvent.ID)
		require.Equal(t, err, storage.ErrNotFound)
		require.Nil(t, storedEvent)
	})

	t.Run("delete unexisting event fails", func(t *testing.T) {
		s := New()

		err := s.DeleteEvent("1")
		require.Equal(t, err, storage.ErrNotFound)
	})
}

func TestStorageListEvents(t *testing.T) {
	s := New()

	event1 := storage.Event{
		ID:       "1",
		Title:    "Event 1",
		Start:    time.Date(2021, time.April, 10, 0, 0, 0, 0, time.UTC),
		Duration: time.Hour,
	}
	event2 := storage.Event{
		ID:       "2",
		Title:    "Event 2",
		Start:    time.Date(2021, time.April, 21, 0, 0, 0, 0, time.UTC),
		Duration: time.Hour,
	}
	event3 := storage.Event{
		ID:       "3",
		Title:    "Event 3",
		Start:    time.Date(2021, time.May, 10, 0, 0, 0, 0, time.UTC),
		Duration: time.Hour,
	}
	event4 := storage.Event{
		ID:       "4",
		Title:    "Event 4",
		Start:    time.Date(2021, time.April, 11, 0, 0, 0, 0, time.UTC),
		Duration: time.Hour,
	}

	require.NoError(t, s.CreateEvent(event1))
	require.NoError(t, s.CreateEvent(event2))
	require.NoError(t, s.CreateEvent(event3))
	require.NoError(t, s.CreateEvent(event4))

	t.Run("list events for day succeeds", func(t *testing.T) {
		dayEvents, err := s.ListEventsForDay(
			time.Date(2021, time.April, 10, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1}, dayEvents)
	})

	t.Run("list events for week succeeds", func(t *testing.T) {
		weekEvents, err := s.ListEventsForWeek(
			time.Date(2021, time.April, 5, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1, event4}, weekEvents)
	})

	t.Run("list events for month succeeds", func(t *testing.T) {
		monthEvents, err := s.ListEventsForMonth(
			time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1, event4, event2}, monthEvents)
	})
}
