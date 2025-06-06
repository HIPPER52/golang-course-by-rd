package client

import (
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) Register(ctx *fiber.Ctx) error {
	var payload dto.RegisterClientDTO
	if err := ctx.BodyParser(&payload); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request")
	}

	client, err := h.svc.Client.RegisterClient(ctx.Context(), payload)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	dialog := &models.QueuedDialog{
		DialogBase: models.DialogBase{
			ID:          ulid.Make().String(),
			ClientID:    client.ID,
			ClientName:  client.Name,
			ClientPhone: client.Phone,
			ClientIP:    ctx.IP(),
			StartedAt:   client.CreatedAt,
		},
	}

	if err := h.svc.QueuedDialog.Add(ctx.Context(), dialog); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to create dialog")
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"client": client,
		"roomID": dialog.ID,
	})
}
