package storage

import (
	"errors"
	"time"

	"github.com/google/uuid"
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
	title string,
	descrip string,
	sTime time.Time,
	eTime time.Time,
	userId string,
	callDuration time.Duration,
) *Event {
	id := uuid.New()
	return &Event{
		ID:           id.String(),
		Title:        title,
		Description:  descrip,
		StartTime:    sTime,
		EndTime:      eTime,
		UserID:       userId,
		CallDuration: callDuration,
	}
}
