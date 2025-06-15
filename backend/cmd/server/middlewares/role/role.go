package role

import (
	"course_project/internal/constants"
	"course_project/internal/constants/roles"
	"github.com/gofiber/fiber/v2"
)

const ctxRoleKey = constants.CONTEXT_ROLE

type Middleware struct{}

func NewMiddleware() *Middleware {
	return &Middleware{}
}

func (m *Middleware) RequireRoles(allowedRoles ...roles.Role) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userRoleVal := ctx.Locals(ctxRoleKey)
		userRole, ok := userRoleVal.(roles.Role)
		if !ok {
			return fiber.ErrForbidden
		}

		for _, role := range allowedRoles {
			if userRole == role {
				return ctx.Next()
			}
		}

		return fiber.ErrForbidden
	}
}
