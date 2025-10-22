package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"ride-hail/internal/core/domain/models"
	"time"
)

// уровни логов (по желанию)
const (
	LevelDebug = "DEBUG"
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
)

// Logger — основной логгер, создаётся на уровне сервиса
type Logger struct {
	service  string
	hostname string
	slog     *slog.Logger
}

// New создает новый логгер.
// pretty=true → красиво в консоли (dev)
// pretty=false → компактный JSON (prod)
func New(service string, pretty bool) *Logger {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown-host"
	}

	var handler slog.Handler
	if pretty {
		handler = NewPrettyJSONHandler(os.Stdout, slog.LevelDebug)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}

	return &Logger{
		service:  service,
		hostname: hostname,
		slog:     slog.New(handler),
	}
}

// Func создаёт логгер для конкретной функции (или метода)
func (l *Logger) Func(prefix string) *FuncLogger {
	return &FuncLogger{
		service:  l.service,
		prefix:   prefix,
		hostname: l.hostname,
		slog:     l.slog,
	}
}

// FuncLogger — логгер уровня функции
type FuncLogger struct {
	service  string
	prefix   string
	hostname string
	slog     *slog.Logger
}

func (f *FuncLogger) log(ctx context.Context, level slog.Level, action, message string, fields ...interface{}) {
	reqID := GetRequestID(ctx)
	userID := GetUserID(ctx)

	attrs := []slog.Attr{
		slog.String("service", f.service),
		slog.String("func_name", f.prefix),
		slog.String("action", action),
		slog.String("hostname", f.hostname),
	}

	if reqID != "" {
		attrs = append(attrs, slog.String("request_id", reqID))
	}

	if userID != "" {
		attrs = append(attrs, slog.String("user_id", userID))
	}

	for i := 0; i < len(fields)-1; i += 2 {
		key, ok := fields[i].(string)
		if !ok {
			continue
		}
		attrs = append(attrs, slog.Any(key, fields[i+1]))
	}

	f.slog.LogAttrs(ctx, level, message, attrs...)
}

// Методы уровней логирования
func (f *FuncLogger) Debug(ctx context.Context, action, message string, fields ...interface{}) {
	f.log(ctx, slog.LevelDebug, action, message, fields...)
}
func (f *FuncLogger) Info(ctx context.Context, action, message string, fields ...interface{}) {
	f.log(ctx, slog.LevelInfo, action, message, fields...)
}
func (f *FuncLogger) Warn(ctx context.Context, action, message string, fields ...interface{}) {
	f.log(ctx, slog.LevelWarn, action, message, fields...)
}
func (f *FuncLogger) Error(ctx context.Context, action, message string, fields ...interface{}) {
	f.log(ctx, slog.LevelError, action, message, fields...)
}

// контекст для request_id
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

func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, models.GetUserIDKey(), userID)
}

func GetUserID(ctx context.Context) string {
	if v := ctx.Value(models.GetUserIDKey()); v != nil {
		if id, ok := v.(string); ok {
			return id
		}
	}
	return ""
}

func WithRole(ctx context.Context, role string) context.Context {
	return context.WithValue(ctx, models.GetRoleKey(), role)
}

func GetRole(ctx context.Context) string {
	if v := ctx.Value(models.GetRoleKey()); v != nil {
		if r, ok := v.(string); ok {
			return r
		}
	}
	return ""
}

////////////////////////////////////////////////////////////////////////////////
// PRETTY JSON HANDLER — красивый вывод JSON для dev
////////////////////////////////////////////////////////////////////////////////

type PrettyJSONHandler struct {
	out   io.Writer
	level slog.Leveler
}

func NewPrettyJSONHandler(out io.Writer, level slog.Leveler) *PrettyJSONHandler {
	return &PrettyJSONHandler{out: out, level: level}
}

func (h *PrettyJSONHandler) Enabled(_ context.Context, lvl slog.Level) bool {
	return lvl >= h.level.Level()
}

func (h *PrettyJSONHandler) Handle(_ context.Context, r slog.Record) error {
	data := make(map[string]interface{})
	data["time"] = r.Time.Format(time.RFC3339)
	data["level"] = r.Level.String()
	data["msg"] = r.Message

	r.Attrs(func(a slog.Attr) bool {
		data[a.Key] = a.Value.Any()
		return true
	})

	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return err
	}
	_, err := h.out.Write(buf.Bytes())
	return err
}

func (h *PrettyJSONHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

func (h *PrettyJSONHandler) WithGroup(name string) slog.Handler {
	return h
}
