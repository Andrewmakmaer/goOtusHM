package app

import (
	"context"
)

type App struct {
}

type Config interface {
}

type Logger interface {
	Debug(string)
	Info(string)
	Warn(string)
	Error(string)
}

type Storage interface {
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
