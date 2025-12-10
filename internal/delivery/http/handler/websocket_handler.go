package handler

import (
	"bufio"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	manager *WebSocketManager
}

func NewWebSocketHandler(manager *WebSocketManager) *WebSocketHandler {
	return &WebSocketHandler{manager: manager}
}

func (h *WebSocketHandler) Handle(c *gin.Context) {
	userId := c.Query("userId")
	if userId == "" {
		c.JSON(400, gin.H{"error": "userId is required"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	// Add client
	h.manager.AddClient(userId, conn)
	defer h.manager.RemoveClient(userId)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		prompt := url.QueryEscape(string(msg))
		resp, err := http.Get("http://localhost:9090/ask?question=" + prompt)
		if err != nil {
			h.manager.SendToUser(userId, "AI error")
			continue
		}
		defer resp.Body.Close()

		reader := bufio.NewScanner(resp.Body)

		// Stream token → specific user
		for reader.Scan() {
			token := reader.Text()
			h.manager.SendToUser(userId, token)
		}
	}
}
