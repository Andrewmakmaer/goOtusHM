package callcearcher

import (
	"fmt"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
)

type Storage interface {
	GetAllEvents() ([]storage.Event, error)
	UpdateEvent(storage.Event) error
	DeleteEvent(string, string) error
}

type Broker interface {
	SendMessage(storage.Event) error
}

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

func SendCallCandidate(store Storage, broker Broker, logg Logger) error {
	sendCandidate, err := searchCallCandidate(store, logg)
	if err != nil {
		return err
	}

	for _, item := range sendCandidate {
		logg.Info("message", "send notification", "UserID", item.UserID, "title", item.Title)
		broker.SendMessage(item)
	}
	return nil
}

func DeleteOldEvents(store Storage, logg Logger) error {
	delCandidates, err := searchDelCandidate(store)
	if err != nil {
		return err
	}

	for _, item := range delCandidates {
		logg.Info("message", "delete events", "UserID", item.UserID, "title", item.Title)
		store.DeleteEvent(item.ID, item.UserID)
	}
	return nil
}

func searchCallCandidate(store Storage, logg Logger) ([]storage.Event, error) {
	events, err := store.GetAllEvents()
	if err != nil {
		return []storage.Event{}, err
	}

	var resultEvents []storage.Event
	now := time.Now()
	for _, item := range events {
		if item.StartTime.Sub(now) < item.CallDuration && now.Before(item.StartTime) && item.CallDuration > 0*time.Second {
			resultEvents = append(resultEvents, item)
			err := resetingDuration(item, store)
			if err != nil {
				logg.Error("message", err.Error())
			}
		}
	}
	return resultEvents, nil
}

func resetingDuration(item storage.Event, store Storage) error {
	item.CallDuration = 0 * time.Second
	err := store.UpdateEvent(item)
	if err != nil {
		return err
	}
	return nil
}

func searchDelCandidate(store Storage) ([]storage.Event, error) {
	events, err := store.GetAllEvents()
	if err != nil {
		return []storage.Event{}, err
	}

	var resultEvents []storage.Event
	olddestEventsTime := time.Now().AddDate(-1, 0, 0)
	for _, item := range events {
		if item.EndTime.Before(olddestEventsTime) {
			fmt.Println(olddestEventsTime, item.EndTime)
			resultEvents = append(resultEvents, item)
		}
	}
	return resultEvents, nil
}
