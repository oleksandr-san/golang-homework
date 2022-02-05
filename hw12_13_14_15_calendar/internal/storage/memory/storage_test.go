package memorystorage

import (
	"context"
	"testing"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	t.Run("describe existing event succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "Test"}

		err := s.CreateEvent(context.TODO(), testEvent)
		require.NoError(t, err)

		event, err := s.ReadEvent(context.TODO(), "1", "")
		require.NoError(t, err)

		require.Equal(t, testEvent, *event)
	})

	t.Run("describe unexisting event fails", func(t *testing.T) {
		s := New()
		event, err := s.ReadEvent(context.TODO(), "1", "")

		require.Nil(t, event)
		require.Equal(t, err, storage.ErrNotFound)
	})

	t.Run("update existing even succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "A"}
		err := s.CreateEvent(context.TODO(), testEvent)
		require.NoError(t, err)

		updatedEvent := storage.Event{ID: "1", Title: "B"}
		err = s.UpdateEvent(context.TODO(), updatedEvent)
		require.NoError(t, err)

		storedEvent, err := s.ReadEvent(context.TODO(), "1", "")
		require.NoError(t, err)
		require.NotNil(t, storedEvent)
		require.Equal(t, updatedEvent, *storedEvent)
	})

	t.Run("update unexising event fails", func(t *testing.T) {
		s := New()
		err := s.UpdateEvent(context.TODO(), storage.Event{ID: "1"})
		require.Equal(t, err, storage.ErrNotFound)
	})

	t.Run("delete existing event succeeds", func(t *testing.T) {
		s := New()

		testEvent := storage.Event{ID: "1", Title: "A"}
		err := s.CreateEvent(context.TODO(), testEvent)
		require.NoError(t, err)

		err = s.DeleteEvent(context.TODO(), testEvent.ID, "")
		require.NoError(t, err)

		storedEvent, err := s.ReadEvent(context.TODO(), testEvent.ID, "")
		require.Equal(t, err, storage.ErrNotFound)
		require.Nil(t, storedEvent)
	})

	t.Run("delete unexisting event fails", func(t *testing.T) {
		s := New()

		err := s.DeleteEvent(context.TODO(), "1", "")
		require.Equal(t, err, storage.ErrNotFound)
	})
}

func TestStorageListEvents(t *testing.T) {
	s := New()

	event1 := storage.Event{
		ID:        "1",
		Title:     "Event 1",
		StartTime: time.Date(2021, time.April, 10, 0, 0, 0, 0, time.UTC),
	}
	event2 := storage.Event{
		ID:        "2",
		Title:     "Event 2",
		StartTime: time.Date(2021, time.April, 21, 0, 0, 0, 0, time.UTC),
	}
	event3 := storage.Event{
		ID:        "3",
		Title:     "Event 3",
		StartTime: time.Date(2021, time.May, 10, 0, 0, 0, 0, time.UTC),
	}
	event4 := storage.Event{
		ID:        "4",
		Title:     "Event 4",
		StartTime: time.Date(2021, time.April, 11, 0, 0, 0, 0, time.UTC),
	}

	require.NoError(t, s.CreateEvent(context.TODO(), event1))
	require.NoError(t, s.CreateEvent(context.TODO(), event2))
	require.NoError(t, s.CreateEvent(context.TODO(), event3))
	require.NoError(t, s.CreateEvent(context.TODO(), event4))

	t.Run("list events for day succeeds", func(t *testing.T) {
		dayEvents, err := s.ListEventsForDay(
			context.TODO(),
			"",
			time.Date(2021, time.April, 10, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1}, dayEvents)
	})

	t.Run("list events for week succeeds", func(t *testing.T) {
		weekEvents, err := s.ListEventsForWeek(
			context.TODO(),
			"",
			time.Date(2021, time.April, 5, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1, event4}, weekEvents)
	})

	t.Run("list events for month succeeds", func(t *testing.T) {
		monthEvents, err := s.ListEventsForMonth(
			context.TODO(),
			"",
			time.Date(2021, time.April, 1, 0, 0, 0, 0, time.UTC))
		require.NoError(t, err)

		require.Equal(t, []storage.Event{event1, event4, event2}, monthEvents)
	})
}
