package logger

import (
	"course_project/internal/services/logger"
	"github.com/gofiber/fiber/v2"
	"log/slog"
)

func getLoggerAttrs(ctx *fiber.Ctx) []slog.Attr {
	av := ctx.Locals(logger.CtxValueKey{})
	if av == nil {
		return []slog.Attr{}
	}

	res, ok := av.([]slog.Attr)
	if !ok {
		return []slog.Attr{}
	}

	return res
}

func mergeAttrs(ctx *fiber.Ctx, attrs []slog.Attr) []slog.Attr {
	existing := getLoggerAttrs(ctx)
	return append(existing, attrs...)
}

func SetLoggerAttrs(ctx *fiber.Ctx, attrs ...slog.Attr) {
	attrs = mergeAttrs(ctx, attrs)
	ctx.Locals(logger.CtxValueKey{}, attrs)
}
