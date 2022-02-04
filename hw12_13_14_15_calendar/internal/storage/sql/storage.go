package sqlstorage

import (
	"context"
	"time"

	"github.com/alexandera5/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Logger interface {
	Debugf(fmt string, v ...interface{})
	Errorf(fmt string, v ...interface{})
}

type Storage struct {
	driverName     string
	dataSourceName string
	db             *sqlx.DB
	logger         Logger
}

func New(driverName, dataSourceName string, logger Logger) *Storage {
	return &Storage{
		driverName:     driverName,
		dataSourceName: dataSourceName,
		logger:         logger,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, s.driverName, s.dataSourceName)
	if err != nil {
		return err
	}

	s.db = db
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.db.Close()
}

func (s *Storage) CreateEvent(ctx context.Context, event storage.Event) error {
	query := `INSERT INTO events(
		id, owner, title, descr,
		start_date, start_time, end_date, end_time
	) values ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := s.db.ExecContext(
		ctx, query,
		event.ID, event.OwnerID, event.Title, event.Description,
		event.StartTime.Format("2019-12-31"),
		event.StartTime.Format("23:00:00"),
		event.EndTime.Format("2019-12-31"),
		event.EndTime.Format("23:00:00"),
	)
	if err != nil {
		return err
	}

	s.logger.Debugf("event created")
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	query := `UPDATE events set
		title = $1, descr = $2,
	) values ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := s.db.ExecContext(
		ctx, query,
		event.ID, event.OwnerID, event.Title, event.Description,
		event.StartTime.Format("2019-12-31"),
		event.StartTime.Format("23:00:00"),
		event.EndTime.Format("2019-12-31"),
		event.EndTime.Format("23:00:00"),
	)
	if err != nil {
		return err
	}

	s.logger.Debugf("event updated")
	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID, ownerID string) error {
	query := "DELETE FROM events WHERE id = $1 AND owner = $2"

	_, err := s.db.ExecContext(ctx, query, eventID, ownerID)
	if err != nil {
		return err
	}

	s.logger.Debugf("event deleted")
	return nil
}

func parseDateTime(dateVal, timeVal string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", dateVal+" "+timeVal)
}

func (s *Storage) queryEvents(ctx context.Context, query string, args ...interface{}) ([]storage.Event, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []storage.Event

	for rows.Next() {
		var ID, ownerID, title, descr string
		var startDateRaw, startTimeRaw, endDateRaw, endTimeRaw string

		err := rows.Scan(
			&ID, &ownerID, &title, &descr,
			&startDateRaw, &startTimeRaw,
			&endDateRaw, &endTimeRaw)
		if err != nil {
			s.logger.Errorf("queryEvents: %v", err)
			return nil, err
		}

		startTime, err := parseDateTime(startDateRaw, startTimeRaw)
		if err != nil {
			return nil, err
		}

		endTime, err := parseDateTime(endDateRaw, endTimeRaw)
		if err != nil {
			return nil, err
		}

		events = append(events, storage.Event{
			ID:          ID,
			OwnerID:     ownerID,
			Title:       title,
			Description: descr,
			StartTime:   startTime,
			EndTime:     endTime,
		})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

func (s *Storage) ReadEvent(ctx context.Context, eventID, ownerID string) (*storage.Event, error) {
	query := `SELECT id, owner, title, descr, start_date, start_time, end_date, end_time
	FROM events 
	WHERE id = $1 AND owner = $2;
	`

	events, err := s.queryEvents(ctx, query, eventID, ownerID)
	if err != nil {
		return nil, err
	}

	if len(events) == 1 {
		return nil, storage.ErrNotFound
	}

	return &events[0], nil
}

func (s *Storage) ListEventsForDay(ctx context.Context, ownerID string, day time.Time) ([]storage.Event, error) {
	query := `SELECT id, owner, title, descr, start_date, start_time, end_date, end_time
	FROM events 
	WHERE owner = $1 AND start_date = $2
	ORDER BY start_time;
	`

	return s.queryEvents(ctx, query, ownerID, day.Format("2019-12-31"))
}

func (s *Storage) ListEventsForWeek(ctx context.Context, ownerID string, firstDay time.Time) ([]storage.Event, error) {
	query := `SELECT id, owner, title, descr, start_date, start_time, end_date, end_time
	FROM events 
	WHERE owner = $1
	AND date_part("year", start_date) = $2 AND date_part("week", start_date) = $3
	ORDER BY start_time;
	`

	year, week := firstDay.ISOWeek()
	return s.queryEvents(ctx, query, ownerID, year, week)
}

func (s *Storage) ListEventsForMonth(ctx context.Context, ownerID string, firstDay time.Time) ([]storage.Event, error) {
	query := `SELECT id, owner, title, descr, start_date, start_time, end_date, end_time
	FROM events 
	WHERE owner = $1
	AND date_part("year", start_date) = $2 AND date_part("month", start_date) = $3
	ORDER BY start_time;
	`

	return s.queryEvents(ctx, query, ownerID, firstDay.Year(), firstDay.Month())
}
