package ws

import (
	"context"
	"course_project/internal/constants"
	"course_project/internal/constants/consumer"
	"course_project/internal/constants/message"
	"course_project/internal/constants/roles"
	wsevent "course_project/internal/constants/ws"
	"course_project/internal/models"
	"course_project/internal/services"
	"course_project/internal/services/logger"
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"github.com/oklog/ulid/v2"
	"time"
)

var allowedEventsByRole = map[roles.Role]map[string]struct{}{
	roles.Client: {
		wsevent.Init:    {},
		wsevent.Message: {},
	},
	roles.Operator: {
		wsevent.Message:      {},
		wsevent.DialogTaken:  {},
		wsevent.DialogClosed: {},
	},
	roles.Admin: {
		wsevent.Message:      {},
		wsevent.DialogTaken:  {},
		wsevent.DialogClosed: {},
	},
}

type ChatGateway struct {
	roomManager *RoomManager
	svc         *services.Services
}

func NewChatGateway(rm *RoomManager, svc *services.Services) *ChatGateway {
	return &ChatGateway{roomManager: rm, svc: svc}
}

type EventPayload struct {
	Event string          `json:"event"`
	Data  json.RawMessage `json:"data"`
}

func (g *ChatGateway) HandleConnection(conn *websocket.Conn) {
	defer conn.Close()

	role := roles.Client
	senderID := ""

	if roleRaw := conn.Locals(constants.CONTEXT_ROLE); roleRaw != nil {
		if parsed, ok := roleRaw.(roles.Role); ok {
			role = parsed
		}
	}

	if userIDAny := conn.Locals(constants.CONTEXT_USER_ID); userIDAny != nil {
		if idPtr, ok := userIDAny.(*string); ok && idPtr != nil {
			senderID = *idPtr
		}
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err,
				websocket.CloseNormalClosure,
				websocket.CloseGoingAway,
			) {
				return
			}

			logger.Error(context.Background(), fmt.Errorf("read error: %v", err))
			return
		}

		var payload EventPayload
		if err := json.Unmarshal(msg, &payload); err != nil {
			logger.Error(context.Background(), fmt.Errorf("invalid format: %v", err))
			_ = conn.WriteMessage(websocket.TextMessage, []byte("invalid format"))
			continue
		}

		if !isEventAllowed(role, payload.Event) {
			_ = conn.WriteMessage(websocket.TextMessage, []byte("unauthorized event"))
			continue
		}

		switch payload.Event {
		case wsevent.Init:
			senderID = g.handleInit(conn, payload.Data)
		case wsevent.Message:
			g.handleMessage(conn, payload.Data, senderID)
		case wsevent.DialogTaken:
			g.handleDialogTaken(conn, payload.Data)
		case wsevent.DialogClosed:
			g.handleDialogClosed(conn, payload.Data)

		default:
			_ = conn.WriteMessage(websocket.TextMessage, []byte("unknown event"))
		}
	}
}

func isEventAllowed(role roles.Role, event string) bool {
	events, ok := allowedEventsByRole[role]
	if !ok {
		return false
	}
	_, exists := events[event]
	return exists
}

func (g *ChatGateway) handleInit(conn *websocket.Conn, data json.RawMessage) string {
	var payload struct {
		RoomID   string `json:"room_id"`
		ClientID string `json:"client_id"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Error(context.Background(), fmt.Errorf("init unmarshal error: %v", err))
		return ""
	}

	g.roomManager.JoinRoom(payload.RoomID, conn)
	logger.Info(context.Background(), fmt.Sprintf("Client %s joined room %s", payload.ClientID, payload.RoomID))

	return payload.ClientID
}

func (g *ChatGateway) handleMessage(conn *websocket.Conn, data json.RawMessage, senderID string) {
	var payload struct {
		RoomID string `json:"room_id"`
		Text   string `json:"text"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Error(context.Background(), fmt.Errorf("message unmarshal error: %v", err))
		return
	}

	msg := models.Message{
		ID:       ulid.Make().String(),
		RoomID:   payload.RoomID,
		SenderID: senderID,
		Content:  payload.Text,
		SentAt:   time.Now().UTC(),
		Type:     message.Text,
	}

	if err := g.svc.Producer.Publish(consumer.TypeSaveMessage, msg); err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to publish message to RabbitMQ: %v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("error publishing message"))
		return
	}

	out := map[string]interface{}{
		"event": wsevent.Message,
		"data": map[string]string{
			"room_id":   payload.RoomID,
			"text":      payload.Text,
			"sender_id": senderID,
		},
	}

	bytes, _ := json.Marshal(out)
	g.roomManager.BroadcastMessage(payload.RoomID, websocket.TextMessage, bytes)

	logger.Info(context.Background(), fmt.Sprintf("%s send message to: %s", senderID, payload.RoomID))
}

func (g *ChatGateway) handleDialogTaken(conn *websocket.Conn, data json.RawMessage) {
	var payload struct {
		RoomID string `json:"room_id"`
		Info   string `json:"info"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Error(context.Background(), fmt.Errorf("dialog_taken unmarshal error: %v", err))
		return
	}

	userIDAny := conn.Locals(constants.CONTEXT_USER_ID)
	opIDPtr, ok := userIDAny.(*string)
	if !ok || opIDPtr == nil {
		logger.Error(context.Background(), fmt.Errorf("missing operator ID in ws context"))
		_ = conn.WriteMessage(websocket.CloseMessage, []byte("unauthorized"))
		return
	}
	operatorID := *opIDPtr

	if err := g.svc.Mover.TakeDialog(context.Background(), payload.RoomID, operatorID); err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to take dialog %s: %v", payload.RoomID, err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("error taking dialog"))
		return
	}

	dialog, err := g.svc.ActiveDialog.FindByID(context.Background(), payload.RoomID)
	if err != nil {
		logger.Error(context.Background(), fmt.Errorf("cannot find taken dialog: %v", err))
		return
	}

	g.roomManager.JoinRoom(dialog.ID, conn)
	logger.Info(context.Background(), fmt.Sprintf("Operator %s joined to room: %s", operatorID, dialog.ID))

	out := map[string]interface{}{
		"event": wsevent.DialogTaken,
		"data": map[string]string{
			"room_id":     dialog.ID,
			"client_name": dialog.ClientName,
		},
	}

	bytes, _ := json.Marshal(out)
	g.roomManager.BroadcastMessage(payload.RoomID, websocket.TextMessage, bytes)
	g.roomManager.BroadcastMessage(wsevent.RoomOperators, websocket.TextMessage, bytes)
}

func (g *ChatGateway) handleDialogClosed(conn *websocket.Conn, data json.RawMessage) {
	var payload struct {
		RoomID string `json:"room_id"`
		Info   string `json:"info"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		logger.Error(context.Background(), fmt.Errorf("dialog_closed unmarshal error: %v", err))
		return
	}

	if err := g.svc.Mover.CloseDialog(context.Background(), payload.RoomID); err != nil {
		logger.Error(context.Background(), fmt.Errorf("failed to close dialog: %v", err))
		_ = conn.WriteMessage(websocket.TextMessage, []byte("error closing dialog"))
		return
	}

	msg := map[string]interface{}{
		"event": wsevent.DialogClosed,
		"data": map[string]string{
			"room_id": payload.RoomID,
			"info":    payload.Info,
		},
	}
	msgBytes, _ := json.Marshal(msg)

	g.roomManager.BroadcastMessage(payload.RoomID, websocket.TextMessage, msgBytes)

	g.roomManager.BroadcastMessage(wsevent.RoomOperators, websocket.TextMessage, msgBytes)
}
