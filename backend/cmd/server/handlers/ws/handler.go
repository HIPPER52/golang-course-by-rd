package ws

import (
	"course_project/internal/services"
)

type Handler struct {
	Manager         *ConnectionManager
	Rooms           *RoomManager
	ClientHandler   *ClientHandler
	OperatorHandler *OperatorHandler
	ChatGateway     *ChatGateway
}

func NewHandler(svcs *services.Services) *Handler {
	manager := NewConnectionManager()
	rooms := NewRoomManager()
	chatGateway := NewChatGateway(rooms, svcs)

	return &Handler{
		Manager:         manager,
		Rooms:           rooms,
		ClientHandler:   NewClientHandler(svcs, manager, rooms, chatGateway),
		OperatorHandler: NewOperatorHandler(svcs, manager, rooms, chatGateway),
	}
}
