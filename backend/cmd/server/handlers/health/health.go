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

// Health godoc
// @Summary      Health check
// @Description  Returns status of the service
// @Tags         health
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse
// @Router       /health [get]
func (h *Handler) Health(ctx *fiber.Ctx) error {
	return ctx.JSON(&HealthResponse{
		Status: "OK",
		Env:    h.cfg.Env,
	})
}
