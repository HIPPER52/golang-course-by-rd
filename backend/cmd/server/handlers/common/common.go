package common

import (
	"course_project/internal/dto"
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	svc *services.Services
}

func NewHandler(svc *services.Services) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) GetDialogMessages(ctx *fiber.Ctx) error {
	var params dto.GetDialogMessagesDTO
	if err := ctx.ParamsParser(&params); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid room_id")
	}

	messages, err := h.svc.Message.FindByRoomID(ctx.Context(), params.RoomID, "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return ctx.JSON(messages)
}
