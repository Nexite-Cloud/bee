package bee

import (
	"context"
	"log/slog"
)

// Logger is an interface for logging in the hive
type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type slogLogger struct {
	logger *slog.Logger
}

// NewSlogLogger creates a new slog logger, if logger is nil, use slog.Default()
func NewSlogLogger(logger *slog.Logger) Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return &slogLogger{logger: logger}
}

// Info logs an info message
func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

// Error logs an error message
func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

// NoLog is a no-op logger
type NoLog struct{}

func (n NoLog) Info(ctx context.Context, msg string, args ...any)  {}
func (n NoLog) Error(ctx context.Context, msg string, args ...any) {}
