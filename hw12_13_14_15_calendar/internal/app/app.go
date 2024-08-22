package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Andrewmakmaer/goOtusHM/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	logg  Logger
	store Storage
}

type Logger interface {
	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
}

type Storage interface {
	AddEvent(e storage.Event) error
	UpdateEvent(updatedEvent storage.Event) error
	DeleteEvent(eventID, userID string) error
	ListEventsDay(userID string, currentDate time.Time) ([]storage.Event, error)
	ListEventsWeek(userID string, currentDate time.Time) ([]storage.Event, error)
	ListEventsMonth(userID string, currentDate time.Time) ([]storage.Event, error)
	GetEvent(eventID, userID string) (storage.Event, error)
}

func New(logger Logger, storage Storage, port, host string) *App {
	return &App{
		logg:  logger,
		store: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, id, userID, title, descrip, bTime, eTime, callDur string) error {
	endTime, err := time.Parse(time.RFC3339, eTime)
	if err != nil {
		a.logg.Error("message", err.Error())
		return err
	}

	beginTime, err := time.Parse(time.RFC3339, bTime)
	if err != nil {
		return err
	}

	callDuration, err := time.ParseDuration(callDur)
	if err != nil {
		return err
	}

	newEvent := storage.NewEvent(id, title, descrip, beginTime, endTime, userID, callDuration)

	err = a.store.AddEvent(*newEvent)
	if err != nil {
		a.logg.Error("failed add event ", err.Error())
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id, userID, title, descrip, bTime, eTime, callDur string) error {
	updatedEvent, err := a.store.GetEvent(id, userID)
	if err != nil {
		return err
	}

	if descrip != "" {
		updatedEvent.Title = descrip
	}

	if title != "" {
		updatedEvent.Title = title
	}

	if bTime != "" {
		updatedEvent.StartTime, err = time.Parse(time.RFC3339, bTime)
		if err != nil {
			return err
		}
	}

	if eTime != "" {
		updatedEvent.EndTime, err = time.Parse(time.RFC3339, eTime)
		if err != nil {
			return err
		}
	}

	if callDur != "" {
		updatedEvent.CallDuration, err = time.ParseDuration(callDur)
		if err != nil {
			return err
		}
	}

	err = a.store.UpdateEvent(updatedEvent)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteEvent(userID, eventID string) error {
	err := a.store.DeleteEvent(eventID, userID)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ListEventDay(userID string) (string, error) {
	events, err := a.store.ListEventsDay(userID, time.Now())
	if err != nil {
		return "", err
	}

	ev, _ := json.Marshal(events)
	return string(ev), nil
}

func (a *App) ListEventWeek(userID string) (string, error) {
	events, err := a.store.ListEventsWeek(userID, time.Now())
	if err != nil {
		return "", err
	}

	ev, _ := json.Marshal(events)
	return string(ev), nil
}

func (a *App) ListEventsMonth(userID string) (string, error) {
	events, err := a.store.ListEventsMonth(userID, time.Now())
	if err != nil {
		return "", err
	}

	ev, _ := json.Marshal(events)
	return string(ev), nil
}
