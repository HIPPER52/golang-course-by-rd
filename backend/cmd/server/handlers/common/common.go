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

// GetDialogMessages godoc
// @Summary      Get messages by room ID
// @Description  Returns all messages associated with a specific room ID
// @Tags         common
// @Accept       json
// @Produce      json
// @Param        room_id  path      string  true  "Room ID"
// @Success      200      {array}   models.Message
// @Failure      400      {string}  string  "Invalid room_id"
// @Failure      500      {string}  string  "Internal server error"
// @Router       /common/messages/{room_id} [get]
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
