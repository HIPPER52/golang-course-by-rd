package auth

import (
	"course_project/cmd/server/middlewares/logger"
	"course_project/internal/constants"
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
	"log/slog"
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

	userID, err := m.svc.Auth.VerifyAuthToken(token)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	usr, err := m.svc.Operator.GetOperatorByID(ctx.Context(), *userID)
	if err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	ctx.Locals(ctxUserIdKey, userID)
	ctx.Locals(ctxRoleKey, usr.Role)
	logger.SetLoggerAttrs(ctx, slog.String(ctxUserIdKey, *userID))
	return ctx.Next()
}
