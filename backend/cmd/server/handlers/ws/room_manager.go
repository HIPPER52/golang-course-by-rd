package ws

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type RoomManager struct {
	rooms map[string]map[*websocket.Conn]bool
	lock  sync.RWMutex
}

func NewRoomManager() *RoomManager {
	return &RoomManager{
		rooms: make(map[string]map[*websocket.Conn]bool),
	}
}

func (rm *RoomManager) JoinRoom(roomID string, conn *websocket.Conn) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if _, ok := rm.rooms[roomID]; !ok {
		rm.rooms[roomID] = make(map[*websocket.Conn]bool)
	}
	rm.rooms[roomID][conn] = true
}

func (rm *RoomManager) LeaveRoom(roomID string, conn *websocket.Conn) {
	rm.lock.Lock()
	defer rm.lock.Unlock()

	if conns, ok := rm.rooms[roomID]; ok {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(rm.rooms, roomID)
		}
	}
}

func (rm *RoomManager) Broadcast(roomID string, message string) {
	rm.BroadcastMessage(roomID, websocket.TextMessage, []byte(message))
}

func (rm *RoomManager) BroadcastMessage(roomID string, messageType int, message []byte) {
	rm.lock.RLock()
	defer rm.lock.RUnlock()

	if conns, ok := rm.rooms[roomID]; ok {
		for conn := range conns {
			_ = conn.WriteMessage(messageType, message)
		}
	}
}
