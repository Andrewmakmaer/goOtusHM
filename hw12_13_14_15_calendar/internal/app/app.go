package app

import (
	"context"
	"fmt"
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
		return err
	}

	beginTime, err := time.Parse(time.RFC3339, eTime)
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
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, id, userID, title, descrip, bTime, eTime, callDur string) error {
	endTime, err := time.Parse(time.RFC3339, eTime)
	if err != nil {
		return err
	}

	beginTime, err := time.Parse(time.RFC3339, eTime)
	if err != nil {
		return err
	}

	callDuration, err := time.ParseDuration(callDur)
	if err != nil {
		return err
	}

	newEvent := storage.NewEvent(id, title, descrip, beginTime, endTime, userID, callDuration)

	err = a.store.UpdateEvent(*newEvent)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) DeleteEvent(userID, eventID string) error {
	err := a.store.DeleteEvent(userID, eventID)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) ListEventDay(userID string) error {
	events, err := a.store.ListEventsDay(userID, time.Now())
	if err != nil {
		return err
	}

	// пока так, потом переделаем на нормально
	for _, item := range events {
		fmt.Println(item)
	}

	return nil
}

func (a *App) ListEventWeek(userID string) error {
	events, err := a.store.ListEventsWeek(userID, time.Now())
	if err != nil {
		return err
	}

	for _, item := range events {
		fmt.Println(item)
	}

	return nil
}

func (a *App) ListEventsMonth(userID string) error {
	events, err := a.store.ListEventsMonth(userID, time.Now())
	if err != nil {
		return err
	}

	for _, item := range events {
		fmt.Println(item)
	}

	return nil
}
