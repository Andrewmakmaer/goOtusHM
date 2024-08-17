package storage

import (
	"errors"
	"time"
)

var (
	ErrDateBusy  = errors.New("date is Busy")
	ErrTimeShift = errors.New("end time of event before event start")
	ErrNoFind    = errors.New("no find event by id")
	ErrNoUser    = errors.New("not found user by id")
)

type Event struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	StartTime    time.Time     `json:"starttime"`
	EndTime      time.Time     `json:"endtime"`
	UserID       string        `json:"userid"`
	CallDuration time.Duration `json:"callduration"`
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
