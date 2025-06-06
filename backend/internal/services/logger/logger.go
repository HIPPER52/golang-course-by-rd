package logger

import (
	"context"
	"log/slog"
	"os"
)

type CtxValueKey struct{}

func init() {
	h := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	l := slog.New(h)
	slog.SetDefault(l)
}

func Info(ctx context.Context, msg string, attrs ...slog.Attr) {
	args := getArgs(mergeAttrs(ctx, attrs))
	slog.Default().InfoContext(ctx, msg, args...)
}

func Error(ctx context.Context, err error, attrs ...slog.Attr) {
	args := getArgs(mergeAttrs(ctx, attrs))
	slog.Default().ErrorContext(ctx, err.Error(), args...)
}

func Panic(ctx context.Context, err error, attrs ...slog.Attr) {
	Error(ctx, err, attrs...)
	panic(err)
}

func Fatal(ctx context.Context, err error, attrs ...slog.Attr) {
	Error(ctx, err, attrs...)
	os.Exit(1)
}

func getArgs(attrs []slog.Attr) []any {
	args := make([]any, len(attrs))
	for i, attr := range attrs {
		args[i] = attr
	}

	return args
}

func getAttrs(ctx context.Context) []slog.Attr {
	av := ctx.Value(CtxValueKey{})
	if av == nil {
		return []slog.Attr{}
	}

	res, ok := av.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}

	return res
}

func mergeAttrs(ctx context.Context, attrs []slog.Attr) []slog.Attr {
	existing := getAttrs(ctx)
	return append(existing, attrs...)
}
