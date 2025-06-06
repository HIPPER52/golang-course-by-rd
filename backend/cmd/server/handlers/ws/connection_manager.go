package ws

import (
	"github.com/gofiber/websocket/v2"
	"sync"
)

type Connection struct {
	UserID string
	Role   string
	Conn   *websocket.Conn
}

type ConnectionManager struct {
	connections map[string]*Connection
	lock        sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*Connection),
	}
}

func (h *ConnectionManager) AddConnection(userID, role string, conn *websocket.Conn) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.connections[userID] = &Connection{UserID: userID, Role: role, Conn: conn}
}

func (h *ConnectionManager) RemoveConnection(userID string) {
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.connections, userID)
}

func (h *ConnectionManager) GetConnection(userID string) (*Connection, bool) {
	h.lock.RLock()
	defer h.lock.RUnlock()
	conn, ok := h.connections[userID]
	return conn, ok
}
