package ws

import (
	"course_project/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

type ClientHandler struct {
	svcs        *services.Services
	mgr         *ConnectionManager
	rooms       *RoomManager
	chatGateway *ChatGateway
}

func NewClientHandler(svcs *services.Services, mgr *ConnectionManager, rooms *RoomManager, chatGateway *ChatGateway) *ClientHandler {
	return &ClientHandler{
		svcs:        svcs,
		mgr:         mgr,
		rooms:       rooms,
		chatGateway: chatGateway,
	}
}

func (h *ClientHandler) Handle() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		h.chatGateway.HandleConnection(c)
	})
}
