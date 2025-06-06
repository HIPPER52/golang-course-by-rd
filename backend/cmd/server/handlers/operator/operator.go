package operator

import (
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetQueuedDialogs(ctx *fiber.Ctx) error {
	dialogs, err := h.svc.QueuedDialog.ListAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load queued dialogs",
		})
	}
	return ctx.JSON(dialogs)
}

func (h *Handler) GetActiveDialogs(ctx *fiber.Ctx) error {
	dialogs, err := h.svc.ActiveDialog.ListAll(ctx.Context())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to load active dialogs",
		})
	}
	return ctx.JSON(dialogs)
}
