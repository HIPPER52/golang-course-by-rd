package logger

import (
	"course_project/internal/services"
	"course_project/internal/services/logger"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"log/slog"
	"time"
)

type Middleware struct {
	svcs *services.Services
}

func NewMiddleware(svcs *services.Services) *Middleware {
	return &Middleware{
		svcs: svcs,
	}
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {
	reqID := ulid.Make().String()
	reqPath := ctx.Path()

	SetLoggerAttrs(ctx, slog.String("req_id", reqID), slog.String("req_path", reqPath))

	logger.Info(ctx.Context(), "request start")

	startedAt := time.Now()
	err := ctx.Next()
	duration := time.Since(startedAt)

	SetLoggerAttrs(ctx, slog.Int64("duration", duration.Milliseconds()))

	fe := &fiber.Error{}
	if errors.As(err, &fe) {
		SetLoggerAttrs(ctx, slog.String("resp_message", fe.Message))
	} else if err != nil {
		SetLoggerAttrs(ctx, slog.String("resp_message", err.Error()))
	}

	SetLoggerAttrs(ctx, slog.Int("resp_status", ctx.Response().StatusCode()))

	logger.Info(ctx.Context(), "request end")
	return err
}
