package health

import (
	"course_project/internal/config"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	cfg *config.Config
}

func NewHandler(cfg *config.Config) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

type HealthResponse struct {
	Status string `json:"status"`
	Env    string `json:"env"`
}

func (h *Handler) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(&HealthResponse{
		Status: "OK",
		Env:    h.cfg.Env,
	})
}
