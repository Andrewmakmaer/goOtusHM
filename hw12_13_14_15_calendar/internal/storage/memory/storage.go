package memorystorage

import (
	"sync"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu    sync.RWMutex
	store map[string][]storage.Event
}

func New() *Storage {
	return &Storage{
		store: make(map[string][]storage.Event),
	}
}

func (s *Storage) AddEvent(e storage.Event) error {
	if e.StartTime.After(e.EndTime) {
		return storage.ErrTimeShift
	}

	s.mu.Lock()
	events := s.store[e.UserID]
	s.mu.Unlock()

	for _, v := range events {
		if e.StartTime.Before(v.EndTime) && v.StartTime.Before(e.EndTime) {
			return storage.ErrDateBusy
		}
	}

	s.mu.Lock()
	s.store[e.UserID] = append(s.store[e.UserID], e)
	s.mu.Unlock()
	return nil
}

func (s *Storage) UpdateEvent(updatedEvent storage.Event) error {
	s.mu.Lock()
	eventList := s.store[updatedEvent.UserID]
	s.mu.Unlock()

	for _, v := range eventList {
		if updatedEvent.StartTime.Before(v.EndTime) && v.StartTime.Before(updatedEvent.EndTime) && v.ID != updatedEvent.ID {
			return storage.ErrDateBusy
		}
	}

	for i := range len(eventList) {
		if eventList[i].ID == updatedEvent.ID {
			eventList[i] = updatedEvent
		}
	}
	s.mu.Lock()
	s.store[updatedEvent.UserID] = eventList
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(eventID, userID string) error {
	s.mu.Lock()
	eventList, ok := s.store[userID]
	s.mu.Unlock()
	if !ok {
		return storage.ErrNoUser
	}

	var i int
	var searchedFlag bool
	for i := range len(eventList) {
		if eventList[i].ID == eventID {
			searchedFlag = true
			break
		}
		searchedFlag = false
	}

	if len(eventList) == 0 || !searchedFlag {
		return storage.ErrNoFind
	}
	eventList[i] = eventList[len(eventList)-1]

	s.mu.Lock()
	s.store[userID] = eventList[:len(eventList)-1]
	s.mu.Unlock()
	return nil
}

func (s *Storage) ListEventsDay(userID string, currentDate time.Time) ([]storage.Event, error) {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userID]
	s.mu.RUnlock()

	dayBegin := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	dayEnd := dayBegin.AddDate(0, 1, 0).Add(-time.Nanosecond)

	for _, v := range events {
		if timeInBetween(dayBegin, dayEnd, currentDate) {
			resultEvents = append(resultEvents, v)
		}
	}
	return resultEvents, nil
}

func (s *Storage) ListEventsWeek(userID string, currentDate time.Time) ([]storage.Event, error) {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userID]
	s.mu.RUnlock()

	year, week := currentDate.ISOWeek()
	weekBegin := time.Date(year, 1, 0, 0, 0, 0, 0, currentDate.Location())
	weekBegin = weekBegin.AddDate(0, 0, (week-1)*7+1)
	for weekBegin.Weekday() != time.Monday {
		weekBegin = weekBegin.AddDate(0, 0, -1)
	}

	weekEnd := weekBegin.Add(7*24*time.Hour - time.Nanosecond)
	for _, v := range events {
		if timeInBetween(weekBegin, weekEnd, currentDate) {
			resultEvents = append(resultEvents, v)
		}
	}
	return resultEvents, nil
}

func (s *Storage) ListEventsMonth(userID string, currentDate time.Time) ([]storage.Event, error) {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userID]
	s.mu.RUnlock()

	monthStart := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	for _, v := range events {
		if timeInBetween(monthStart, monthEnd, currentDate) {
			resultEvents = append(resultEvents, v)
		}
	}
	return resultEvents, nil
}

func (s *Storage) GetEvent(eventID, userID string) (storage.Event, error) {
	s.mu.RLock()
	events := s.store[userID]
	s.mu.RUnlock()

	for _, v := range events {
		if v.ID == eventID {
			return v, nil
		}
	}
	return storage.Event{}, storage.ErrNoFind
}

func timeInBetween(startTime, finishTime, currentTime time.Time) bool {
	if currentTime.Before(finishTime) && currentTime.After(startTime) {
		return true
	}
	return false
}
