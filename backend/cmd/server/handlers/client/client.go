package client

import (
	"course_project/cmd/server/handlers/ws"
	wsevent "course_project/internal/constants/ws"
	"course_project/internal/dto"
	"course_project/internal/models"
	"course_project/internal/services"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/oklog/ulid/v2"
)

type Handler struct {
	svc       *services.Services
	wsHandler *ws.Handler
}

func NewHandler(svc *services.Services, wsHandler *ws.Handler) *Handler {
	return &Handler{svc: svc, wsHandler: wsHandler}
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

	event := map[string]interface{}{
		"event": wsevent.DialogCreated,
		"data": map[string]interface{}{
			"room_id":      dialog.ID,
			"client_name":  dialog.ClientName,
			"client_phone": dialog.ClientPhone,
			"client_ip":    dialog.ClientIP,
		},
	}
	msg, _ := json.Marshal(event)
	h.wsHandler.Rooms.BroadcastMessage(wsevent.RoomOperators, websocket.TextMessage, msg)

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"client": client,
		"roomID": dialog.ID,
	})
}
