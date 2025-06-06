package middlewares

import (
	"course_project/cmd/server/middlewares/auth"
	"course_project/cmd/server/middlewares/logger"
	"course_project/cmd/server/middlewares/role"
	"course_project/internal/services"
)

type Middlewares struct {
	Logger *logger.Middleware
	Auth   *auth.Middleware
	Role   *role.Middleware
}

func NewMiddlewares(svc *services.Services) *Middlewares {
	return &Middlewares{
		Logger: logger.NewMiddleware(svc),
		Auth:   auth.NewMiddleware(svc),
		Role:   role.NewMiddleware(),
	}
}
