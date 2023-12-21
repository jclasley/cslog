package cslog

import (
	"context"
	"log/slog"
)

func FromContext(ctx context.Context, slog *slog.Logger) context.Context {
	return context.WithValue(ctx, slogKey, slog)
}

func FromBackground(slog *slog.Logger) context.Context {
	return FromContext(context.Background(), slog)
}

func logger(ctx context.Context) *slog.Logger {
	l, ok := ctx.Value(slogKey).(*slog.Logger)
	if !ok || l == nil {
		return slog.Default()
	}

	return l
}

type slogKeyT struct{}

var slogKey slogKeyT

func WithAttrs(ctx context.Context, attr ...slog.Attr) context.Context {
	l := buildLogger(ctx, attr...)
	return FromContext(ctx, l)
}

func WithGroup(ctx context.Context, g string) context.Context {
	l := logger(ctx)
	l = l.WithGroup(g)
	return FromContext(ctx, l)
}

func Debug(ctx context.Context, msg string, args ...slog.Attr) {
	l := buildLogger(ctx, args...)
	l.DebugContext(ctx, msg)
}

func Info(ctx context.Context, msg string, args ...slog.Attr) {
	l := buildLogger(ctx, args...)
	l.InfoContext(ctx, msg)
}

func Warn(ctx context.Context, msg string, args ...slog.Attr) {
	l := buildLogger(ctx, args...)
	l.WarnContext(ctx, msg)
}

func Error(ctx context.Context, msg string, args ...slog.Attr) {
	l := buildLogger(ctx, args...)
	l.ErrorContext(ctx, msg)
}

func buildLogger(ctx context.Context, args ...slog.Attr) *slog.Logger {
	l := logger(ctx)
	for _, arg := range args {
		l = l.With(arg)
	}
	return l
}
