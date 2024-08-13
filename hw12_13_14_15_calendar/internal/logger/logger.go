package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type message struct {
	Timestamp time.Time              `json:"timestamp"`
	Level     string                 `json:"level"`
	Data      map[string]interface{} `json:"data"`
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

func (l Logger) log(level string, keyvals ...interface{}) {
	if l.level > loggerMap[level] {
		return
	}
	l.printer(level, keyvals...)
}

func (l Logger) Debug(keyvals ...interface{})    { l.log("DEBUG", keyvals...) }
func (l Logger) Info(keyvals ...interface{})     { l.log("INFO", keyvals...) }
func (l Logger) Warn(keyvals ...interface{})     { l.log("WARN", keyvals...) }
func (l Logger) Error(keyvals ...interface{})    { l.log("ERROR", keyvals...) }
func (l Logger) Critical(keyvals ...interface{}) { l.log("CRITICAL", keyvals...) }

func (l Logger) printer(msgLvl string, keyvals ...interface{}) {
	data := make(map[string]interface{})
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			key, ok := keyvals[i].(string)
			if !ok {
				key = fmt.Sprintf("%v", keyvals[i])
			}
			data[key] = keyvals[i+1]
		}
	}

	message := message{Timestamp: time.Now(), Level: msgLvl, Data: data}

	switch l.loggerType {
	case "json":
		b, err := json.Marshal(message)
		if err != nil {
			fmt.Println("Error marshaling JSON:", err)
			return
		}
		os.Stdout.Write(b)
		os.Stdout.WriteString("\n")
	case "unstructured":
		fmt.Printf("%s %s %v\n", time.Now().Format(time.RFC3339), msgLvl, data)
	}
}
