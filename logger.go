package bee

import (
	"context"
	"log/slog"
)

type Logger interface {
	Info(ctx context.Context, msg string, args ...any)
	Error(ctx context.Context, msg string, args ...any)
}

type slogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(logger *slog.Logger) Logger {
	if logger == nil {
		logger = slog.Default()
	}
	return &slogLogger{logger: logger}
}

func (l *slogLogger) Info(ctx context.Context, msg string, args ...any) {
	l.logger.InfoContext(ctx, msg, args...)
}

func (l *slogLogger) Error(ctx context.Context, msg string, args ...any) {
	l.logger.ErrorContext(ctx, msg, args...)
}

type noLog struct{}

func (n noLog) Info(ctx context.Context, msg string, args ...any)  {}
func (n noLog) Error(ctx context.Context, msg string, args ...any) {}
