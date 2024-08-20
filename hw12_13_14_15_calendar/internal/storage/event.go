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

type EventNotification struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartTime string `json:"starttime"`
	UserID    string `json:"userid"`
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

func NewNotification(id, title string, startTime time.Time, userID string) EventNotification {
	date := startTime.String()
	return EventNotification{ID: id, Title: title, StartTime: date, UserID: userID}
}
