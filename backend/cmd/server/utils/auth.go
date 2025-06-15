package utils

import (
	"course_project/internal/constants"
	"course_project/internal/constants/roles"

	"github.com/gofiber/fiber/v2"
)

func GetUserIDAndRole(ctx *fiber.Ctx) (string, roles.Role, error) {
	userIDRaw := ctx.Locals(constants.CONTEXT_USER_ID)
	roleRaw := ctx.Locals(constants.CONTEXT_ROLE)

	userIDPtr, ok1 := userIDRaw.(*string)
	roleVal, ok2 := roleRaw.(roles.Role)

	if !ok1 || !ok2 || userIDPtr == nil {
		return "", "", fiber.ErrUnauthorized
	}

	return *userIDPtr, roleVal, nil
}
