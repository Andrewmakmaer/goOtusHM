package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy  = errors.New("date is Busy")
	ErrTimeShift = errors.New("end time of event before event start")
)

type Event struct {
	ID           string
	Title        string
	Description  string
	StartTime    time.Time
	EndTime      time.Time
	UserID       string
	CallDuration time.Duration
}

func NewEvent(
	id string,
	title string,
	descrip string,
	sTime time.Time,
	eTime time.Time,
	userID string,
	callDuration time.Duration,
) *Event {
	return &Event{
		ID:           id,
		Title:        title,
		Description:  descrip,
		StartTime:    sTime,
		EndTime:      eTime,
		UserID:       userID,
		CallDuration: callDuration,
	}
}
