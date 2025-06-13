package auth

import (
	"course_project/cmd/server/middlewares/logger"
	"course_project/internal/constants"
	"course_project/internal/services"
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

const headerAuthToken = constants.HEADER_AUTH_TOKEN
const ctxUserIdKey = constants.CONTEXT_USER_ID
const ctxRoleKey = constants.CONTEXT_ROLE

type Middleware struct {
	svc *services.Services
}

func NewMiddleware(svc *services.Services) *Middleware {
	return &Middleware{
		svc: svc,
	}
}

func (m *Middleware) Handle(ctx *fiber.Ctx) error {
	token := ctx.Get(headerAuthToken)
	if token == "" {
		token = ctx.Query("token")
	}

	_, err := m.svc.Auth.VerifyAuthToken(token)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	claims, err := m.svc.Auth.VerifyAuthToken(token)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	ctx.Locals(ctxUserIdKey, claims.UserID)
	ctx.Locals(ctxRoleKey, claims.Role)
	logger.SetLoggerAttrs(ctx, slog.String(ctxUserIdKey, claims.UserID))

	return ctx.Next()
}
