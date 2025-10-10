package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime/debug"
	"time"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelError = "ERROR"
)

type LogEntry struct {
	Timestamp string      `json:"timestamp"`
	Level     string      `json:"level"`
	Service   string      `json:"service"`
	Action    string      `json:"action"`
	Message   string      `json:"message"`
	Hostname  string      `json:"hostname"`
	RequestID string      `json:"request_id"`
	RideID    string      `json:"ride_id,omitempty"`
	Error     *ErrorEntry `json:"error,omitempty"`
}

type ErrorEntry struct {
	Msg   string `json:"msg"`
	Stack string `json:"stack"`
}

type Logger struct {
	service  string
	hostname string
}

func New(service string) *Logger {
	hostname, _ := os.Hostname()
	return &Logger{
		service:  service,
		hostname: hostname,
	}
}

func (l *Logger) log(level, action, message, requestID, rideID string, err error) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339), // ISO 8601
		Level:     level,
		Service:   l.service,
		Action:    action,
		Message:   message,
		Hostname:  l.hostname,
		RequestID: requestID,
		RideID:    rideID,
	}

	if level == LevelError && err != nil {
		entry.Error = &ErrorEntry{
			Msg:   err.Error(),
			Stack: string(debug.Stack()),
		}
	}

	data, _ := json.Marshal(entry)
	fmt.Fprintln(os.Stdout, string(data))
}

func (l *Logger) Info(action, msg, requestID, rideID string) {
	l.log(LevelInfo, action, msg, requestID, rideID, nil)
}

func (l *Logger) Debug(action, msg, requestID, rideID string) {
	l.log(LevelDebug, action, msg, requestID, rideID, nil)
}

func (l *Logger) Error(action, msg, requestID, rideID string, err error) {
	l.log(LevelError, action, msg, requestID, rideID, err)
}
