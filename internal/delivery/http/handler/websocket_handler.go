package handler

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type WebSocketHandler struct {
	manager *WebSocketManager
}

func NewWebSocketHandler(manager *WebSocketManager) *WebSocketHandler {
	return &WebSocketHandler{manager: manager}
}

/* =======================
   STREAM SESSION MANAGER
======================= */

type StreamSession struct {
	Cancel context.CancelFunc
}

var streamSessions sync.Map // userId -> *StreamSession

/* =======================
   MESSAGE STRUCT
======================= */

type ClientMessage struct {
	Type       string `json:"type"` // ask | stop
	Text       string `json:"text"`
	IsFirst    bool   `json:"isFirst"`
	IdCategory int    `json:"idCategory"`
	Topic      int    `json:"topic"`
}

/* =======================
   WS HANDLER
======================= */

func (h *WebSocketHandler) Handle(c *gin.Context) {
	userId := c.Query("userId")
	username := c.Query("username")
	role := c.Query("roleName")

	if userId == "" {
		c.JSON(400, gin.H{"error": "userId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	h.manager.AddClient(userId, conn)
	defer h.manager.RemoveClient(userId)

	// kill stream if socket disconnect
	defer h.stopStream(userId)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			h.manager.SendToUser(userId, `{"type":"error","message":"invalid message format"}`)
			continue
		}

		switch clientMsg.Type {

		// ===== STOP STREAM =====
		case "stop":
			h.stopStream(userId)

		// ===== START STREAM =====
		case "ask":
			// pastikan hanya 1 stream aktif
			h.stopStream(userId)

			ctx, cancel := context.WithCancel(context.Background())
			streamSessions.Store(userId, &StreamSession{Cancel: cancel})

			go h.handleStream(ctx, userId, username, role, clientMsg)
		}
	}
}

/* =======================
   STOP STREAM (CENTRAL)
======================= */

func (h *WebSocketHandler) stopStream(userId string) {
	if v, ok := streamSessions.Load(userId); ok {
		v.(*StreamSession).Cancel()
		streamSessions.Delete(userId)
	}
}

/* =======================
   STREAM TO PYTHON
======================= */

func (h *WebSocketHandler) handleStream(
	ctx context.Context,
	userId string,
	username string,
	role string,
	clientMsg ClientMessage,
) {

	payload := map[string]interface{}{
		"question":   clientMsg.Text,
		"isFirst":    clientMsg.IsFirst,
		"idCategory": clientMsg.IdCategory,
		"topic":      clientMsg.Topic,
		"username":   username,
		"role":       role,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return
	}

	req, err := http.NewRequestWithContext(
		ctx,
		"POST",
		// "http://localhost:9091/ask3/stream",
		"http://localhost:9091/ask3/test",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		h.manager.SendToUser(userId, `{"type":"error","message":"python service error"}`)
		return
	}
	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)

	for {
		select {
		case <-ctx.Done():
			// STREAM DIBATALKAN
			return
		default:
		}

		line, err := reader.ReadBytes('\n')
		if err == io.EOF {
			return
		}
		if err != nil {
			return
		}

		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		var data map[string]interface{}
		if err := json.Unmarshal(line, &data); err != nil {
			continue
		}

		switch data["type"] {

		case "token":
			payload, _ := json.Marshal(map[string]interface{}{
				"type":    "chunk",
				"content": data["data"],
			})
			h.manager.SendToUser(userId, string(payload))

		case "end":
			payload, _ := json.Marshal(map[string]interface{}{
				"type":        "done",
				"topic_id":    data["topic"],
				"category_id": data["category"],
			})
			h.manager.SendToUser(userId, string(payload))
			h.stopStream(userId)
			return

		case "error":
			payload, _ := json.Marshal(map[string]interface{}{
				"type":    "error",
				"message": data["message"],
			})
			h.manager.SendToUser(userId, string(payload))
			h.stopStream(userId)
			return
		}
	}
}
