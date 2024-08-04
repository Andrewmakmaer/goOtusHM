package memorystorage

import (
	"sync"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu    sync.RWMutex //nolint:unused
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

func (s *Storage) DeleteEvent(eventId, userId string) error {
	s.mu.Lock()
	eventList := s.store[userId]
	s.mu.Unlock()
	var i int
	for i = range len(eventList) {
		if eventList[i].ID == eventId {
			break
		}
	}
	eventList[i] = eventList[len(eventList)-1]

	s.mu.Lock()
	s.store[userId] = eventList[:len(eventList)-1]
	s.mu.Unlock()
	return nil
}

func (s *Storage) ListEventsDay(userId string, currentDate time.Time) []storage.Event {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userId]
	s.mu.RUnlock()

	dayBegin := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	dayEnd := dayBegin.AddDate(0, 1, 0).Add(-time.Nanosecond)

	for _, v := range events {
		if timeInBetween(dayBegin, dayEnd, currentDate) {
			resultEvents = append(resultEvents, v)
		}
	}
	return resultEvents
}

func (s *Storage) ListEventsWeek(userId string, currentDate time.Time) []storage.Event {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userId]
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
	return resultEvents
}

func (s *Storage) ListEventsMonth(userId string, currentDate time.Time) []storage.Event {
	var resultEvents []storage.Event
	s.mu.RLock()
	events := s.store[userId]
	s.mu.RUnlock()

	monthStart := time.Date(currentDate.Year(), currentDate.Month(), 1, 0, 0, 0, 0, currentDate.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Nanosecond)

	for _, v := range events {
		if timeInBetween(monthStart, monthEnd, currentDate) {
			resultEvents = append(resultEvents, v)
		}
	}
	return resultEvents
}

func timeInBetween(startTime, finishTime, currentTime time.Time) bool {
	if currentTime.Before(finishTime) && currentTime.After(startTime) {
		return true
	}
	return false
}
