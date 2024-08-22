package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type Storage struct {
	connect *sql.DB
	dbName  string
}

func New(host, dbName, user, pass string) (*Storage, error) {
	connStr := fmt.Sprintf("postgres://%v:%v@%v/%v?sslmode=disable", user, pass, host, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return &Storage{}, err
	}

	err = db.Ping()
	if err != nil {
		return &Storage{}, fmt.Errorf("failed connect to database: %w", err)
	}

	return &Storage{connect: db, dbName: dbName}, nil
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.connect.Close()
	if err != nil {
		return err
	}
	return nil
}

// Это на будущее, но я пока не применил.
func (s *Storage) RunMigrations(migrationsDir string) error {
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		return fmt.Errorf("could not get absolute path: %w", err)
	}

	driver, err := postgres.WithInstance(s.connect, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("could not create database driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+absPath,
		s.dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("could not create migrate instance: %w", err)
	}

	if err := m.Up(); err != nil && errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	fmt.Println("Migrations applied successfully")
	return nil
}

func (s *Storage) AddEvent(e storage.Event) error {
	fmt.Println(e)
	if err := s.checkEventOverlap(e); err != nil {
		return err
	}

	query := `INSERT INTO events(id, title, description, starttime, endtime, userID, callduration)
              VALUES ($1, $2, $3, $4, $5, $6, $7)`

	fmt.Println(e.StartTime, e.EndTime)
	_, err := s.connect.Exec(query, e.ID, e.Title, e.Description, e.StartTime,
		e.EndTime, e.UserID, e.CallDuration.String())
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(updatedEvent storage.Event) error {
	if err := s.checkEventOverlap(updatedEvent); err != nil {
		return err
	}

	query := `UPDATE events 
              SET title = $1, description = $2, starttime = $3, endtime = $4, callduration = $5 
              WHERE id = $6 AND userID = $7`

	result, err := s.connect.Exec(query, updatedEvent.Title, updatedEvent.Description,
		updatedEvent.StartTime, updatedEvent.EndTime,
		updatedEvent.CallDuration.String(), updatedEvent.ID, updatedEvent.UserID)
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

func (s *Storage) DeleteEvent(eventID, userID string) error {
	query := `DELETE FROM events WHERE id = $1 AND userID = $2`

	result, err := s.connect.Exec(query, eventID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event with id %v for user %v", eventID, userID)
	}

	return nil
}

func (s *Storage) ListEventsDay(userID string, currentDate time.Time) ([]storage.Event, error) {
	startOfDay := time.Date(currentDate.Year(), currentDate.Month(), currentDate.Day(), 0, 0, 0, 0, currentDate.Location())
	endOfDay := startOfDay.AddDate(0, 0, 1)
	return s.listEvents(userID, startOfDay, endOfDay)
}

func (s *Storage) ListEventsWeek(userID string, currentDate time.Time) ([]storage.Event, error) {
	year, week := currentDate.ISOWeek()
	weekBegin := time.Date(year, 1, 0, 0, 0, 0, 0, currentDate.Location())
	weekBegin = weekBegin.AddDate(0, 0, (week-1)*7+1)
	for weekBegin.Weekday() != time.Monday {
		weekBegin = weekBegin.AddDate(0, 0, -1)
	}

	weekEnd := weekBegin.Add(7*24*time.Hour - time.Nanosecond)
	return s.listEvents(userID, weekBegin, weekEnd)
}

func (s *Storage) ListEventsMonth(userID string, currentDate time.Time) ([]storage.Event, error) {
	startOfMonth := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	endOfMonth := startOfMonth.AddDate(0, 1, 0)
	return s.listEvents(userID, startOfMonth, endOfMonth)
}

func (s *Storage) GetEvent(eventID, userID string) (storage.Event, error) {
	query := `SELECT * FROM events WHERE id = $1 AND userID = $2`

	var e storage.Event
	var callDurationStr string
	row := s.connect.QueryRow(query, eventID, userID)
	err := row.Scan(&e.ID, &e.Title,
		&e.Description, &e.StartTime, &e.EndTime, &e.UserID, &callDurationStr)
	if err != nil {
		return e, fmt.Errorf("failed to get event: %w", err)
	}

	callDuration, err := time.ParseDuration(callDurationStr)
	if err != nil {
		return e, fmt.Errorf("failed to parse call duration: %w", err)
	}
	e.CallDuration = callDuration

	return e, nil
}

func (s *Storage) GetAllEvents() ([]storage.Event, error) {
	query := `SELECT * FROM events`

	var events []storage.Event
	rows, err := s.connect.Query(query)
	if err != nil {
		return events, fmt.Errorf("failed to get event: %w", err)
	}

	for rows.Next() {
		var e storage.Event
		var callDurationStr string
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.UserID, &callDurationStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		callDuration, err := time.ParseDuration(callDurationStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse call duration: %w", err)
		}
		e.CallDuration = callDuration

		events = append(events, e)
	}

	return events, nil
}

func (s *Storage) listEvents(userID string, startDate, endDate time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, description, starttime, endtime, userID, callduration 
              FROM events 
              WHERE userID = $1 AND starttime >= $2 AND starttime < $3
              ORDER BY starttime`

	rows, err := s.connect.Query(query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()

	var events []storage.Event
	for rows.Next() {
		var e storage.Event
		var callDurationStr string
		err := rows.Scan(&e.ID, &e.Title, &e.Description, &e.StartTime, &e.EndTime, &e.UserID, &callDurationStr)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event: %w", err)
		}
		callDuration, err := time.ParseDuration(callDurationStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse call duration: %w", err)
		}
		e.CallDuration = callDuration

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return events, nil
}

func (s *Storage) checkEventOverlap(event storage.Event) error {
	query := `SELECT COUNT(*) FROM events 
              WHERE userID = $1 AND id != $4
              AND ((starttime <= $2 AND endtime > $2) 
                   OR (starttime < $3 AND endtime >= $3) 
                   OR (starttime >= $2 AND endtime <= $3))`

	var count int
	err := s.connect.QueryRow(query, event.UserID, event.StartTime, event.EndTime, event.ID).Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check event overlap: %w", err)
	}

	if count > 0 {
		return errors.New("event overlaps with existing events")
	}

	return nil
}
