package ws

import (
	"context"
	"course_project/cmd/server/utils"
	ws_event "course_project/internal/constants/ws"
	"course_project/internal/services"
	"course_project/internal/services/logger"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type OperatorHandler struct {
	svcs        *services.Services
	mgr         *ConnectionManager
	rooms       *RoomManager
	chatGateway *ChatGateway
}

func NewOperatorHandler(svcs *services.Services, mgr *ConnectionManager, rooms *RoomManager, gateway *ChatGateway) *OperatorHandler {
	return &OperatorHandler{
		svcs:        svcs,
		mgr:         mgr,
		rooms:       rooms,
		chatGateway: gateway,
	}
}

func (h *OperatorHandler) Handle() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		userID, role, err := utils.GetUserIDAndRole(ctx)
		if err != nil {
			return ctx.SendStatus(fiber.StatusUnauthorized)
		}

		return websocket.New(func(conn *websocket.Conn) {
			h.mgr.AddConnection(userID, string(role), conn)
			defer h.mgr.RemoveConnection(userID)

			h.rooms.JoinRoom(ws_event.RoomOperators, conn)

			logger.Info(context.Background(), fmt.Sprintf("Operator %s connected", userID))

			dialogs, err := h.svcs.ActiveDialog.FindByOperatorID(context.Background(), userID)
			if err != nil {
				logger.Error(context.Background(), fmt.Errorf("cannot restore dialogs for operator %s: %v", userID, err))
			} else {
				for _, d := range dialogs {
					h.rooms.JoinRoom(d.ID, conn)
					logger.Info(context.Background(), fmt.Sprintf("Rejoined room %s for operator %s", d.ID, userID))
				}
			}

			h.chatGateway.HandleConnection(conn)
		})(ctx)
	}
}
