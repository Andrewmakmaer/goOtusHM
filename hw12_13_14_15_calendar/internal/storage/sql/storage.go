package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
	_ "github.com/lib/pq"
)

type Storage struct {
	connect *sql.DB
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	connStr := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	err = db.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed connect to database: %w", err)
	}

	s.connect = db

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.connect.Close()
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) AddEvent(ctx context.Context, e storage.Event) error {
	if err := s.checkEventOverlap(ctx, e); err != nil {
		return err
	}

	query := `INSERT INTO events(id, title, description, starttime, endtime, userid, callduration)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := s.connect.ExecContext(ctx, query, e.ID, e.Title, e.Description, e.StartTime, e.EndTime, e.UserID, e.CallDuration)

	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, updatedEvent storage.Event) error {
	if err := s.checkEventOverlap(ctx, updatedEvent); err != nil {
		return err
	}

	query := `UPDATE events 
              SET title = $1, description = $2, starttime = $3, endtime = $4, callduration = $5 
              WHERE id = $6 AND userid = $7`

	result, err := s.connect.ExecContext(ctx, query, updatedEvent.Title, updatedEvent.Description,
		updatedEvent.StartTime, updatedEvent.EndTime,
		updatedEvent.CallDuration, updatedEvent.ID, updatedEvent.UserID)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event with id %v for user %v", updatedEvent.ID, updatedEvent.UserID)
	}

	return nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventId, userId string) error {
	query := `DELETE FROM events WHERE id = $1 AND userid = $2`

	result, err := s.connect.ExecContext(ctx, query, eventId, userId)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event with id %v for user %v", eventId, userId)
	}

	return nil
}

func (s *Storage) ListEventsDay(ctx context.Context, userId string, currentDate time.Time) ([]storage.Event, error) {
	startOfDay := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	return s.listEvents(ctx, userId, startOfDay, endOfDay)
}

func (s *Storage) ListEventsWeek(ctx context.Context, userId string, currentDate time.Time) ([]storage.Event, error) {
	startOfWeek := currentDate.AddDate(0, 0, -int(currentDate.Weekday()+1))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())
	endOfWeek := startOfWeek.AddDate(0, 0, 7)
	return s.listEvents(ctx, userId, startOfWeek, endOfWeek)
}

func (s *Storage) ListEventsMonth(ctx context.Context, userId string, currentDate time.Time) ([]storage.Event, error) {
	startOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	return s.listEvents(ctx, userId, startOfMonth, endOfMonth)
}

func (s *Storage) listEvents(ctx context.Context, userId string, startDate, endDate time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, description, starttime, endtime, userid, callduration 
              FROM events 
              WHERE userid = $1 AND starttime >= $2 AND starttime < $3
              ORDER BY starttime`

	rows, err := s.connect.QueryContext(ctx, query, userId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.UserID, &e.CallDuration)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

func (s *Storage) checkEventOverlap(ctx context.Context, event storage.Event) error {
	query := `SELECT COUNT(*) FROM events 
              WHERE userid = $1 
              AND ((starttime <= $2 AND endtime > $2) 
                   OR (starttime < $3 AND endtime >= $3) 
                   OR (starttime >= $2 AND endtime <= $3))`

	var count int
	err := s.connect.QueryRowContext(ctx, query, event.UserID, event.StartTime, event.EndTime).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check event overlap: %w", err)
	}

	if count > 0 {
		return errors.New("event overlaps with existing events")
	}

	return nil
}
