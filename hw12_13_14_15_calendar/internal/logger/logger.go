package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type message struct {
	Timestamp time.Time
	Message   string
	Level     string
}

type Logger struct {
	level      uint8
	loggerType string
}

var loggerMap = map[string]uint8{
	"DEBUG":    1,
	"INFO":     2,
	"WARN":     3,
	"ERROR":    4,
	"CRITICAL": 5,
}

func New(levelStr, logType string) *Logger {
	return &Logger{level: loggerMap[levelStr], loggerType: logType}
}

func (l Logger) Debug(msg string) {
	if l.level > 1 {
		return
	}
	l.printer(msg, "DEBUG")
}

func (l Logger) Info(msg string) {
	if l.level > 2 {
		return
	}
	l.printer(msg, "INFO")
}

func (l Logger) Warn(msg string) {
	if l.level > 3 {
		return
	}
	l.printer(msg, "WARN")
}

func (l Logger) Error(msg string) {
	if l.level > 4 {
		return
	}
	l.printer(msg, "ERROR")
}

func (l Logger) Critical(msg string) {
	l.printer(msg, "CRITICAL")
}

func (l Logger) printer(msg string, msgLvl string) {
	message := message{Timestamp: time.Now(), Message: msg, Level: msgLvl}
	switch l.loggerType {
	case "json":
		b, err := json.Marshal(message)
		if err != nil {
			fmt.Println("No write", msg)
			return
		}
		os.Stdout.Write(b)
		os.Stdout.WriteString("\n")
	case "unstructed":
		fmt.Println(time.Now(), msg, msgLvl)
	}
}
