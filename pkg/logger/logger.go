package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"ride-hail/internal/core/domain/models"
	"time"
)

const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

type Logger struct {
	service  string
	hostname string
}

func New(service string) *Logger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}
	return &Logger{
		service:  service,
		hostname: hostname,
	}
}

func (l *Logger) Func(funcName string) *FuncLogger {
	return &FuncLogger{
		service:  l.service,
		funcName: funcName,
		hostname: l.hostname,
	}
}

type FuncLogger struct {
	service  string
	funcName string
	hostname string
}

// Основной метод логирования
func (f *FuncLogger) log(level, action, message string, fields ...interface{}) {
	// сначала фиксированные поля
	order := []struct {
		key string
		val interface{}
	}{
		{"level", level},
		{"timestamp", time.Now().Format(time.RFC3339)},
		{"service", f.service},
		{"func", f.funcName},
		{"action", action},
		{"message", message},
		{"hostname", f.hostname},
	}

	// буфер для сборки JSON вручную
	var buf bytes.Buffer
	buf.WriteByte('{')

	for i, kv := range order {
		b, _ := json.Marshal(kv.val)
		fmt.Fprintf(&buf, `"%s":%s`, kv.key, b)
		if i < len(order)-1 {
			buf.WriteByte(',')
		}
	}

	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		val, _ := json.Marshal(fields[i+1])
		buf.WriteByte(',')
		fmt.Fprintf(&buf, `"%s":%s`, key, val)
	}

	buf.WriteByte('}')
	fmt.Fprintln(os.Stdout, buf.String())
}

func (f *FuncLogger) Debug(action, message string, fields ...interface{}) {
	f.log(LevelDebug, action, message, fields...)
}
func (f *FuncLogger) Info(action, message string, fields ...interface{}) {
	f.log(LevelInfo, action, message, fields...)
}
func (f *FuncLogger) Warn(action, message string, fields ...interface{}) {
	f.log(LevelWarn, action, message, fields...)
}
func (f *FuncLogger) Error(action, message string, fields ...interface{}) {
	f.log(LevelError, action, message, fields...)
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, models.GetRequestIDKey(), requestID)
}

func GetRequestID(ctx context.Context) string {
	if v := ctx.Value(models.GetRequestIDKey()); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}
