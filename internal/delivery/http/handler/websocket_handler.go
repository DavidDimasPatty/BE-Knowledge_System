package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"

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
	idCategory := c.Query("idCategory")
	topic := c.Query("topic")
	username := c.Query("username")

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
	type AIResponse struct {
		Answer     string `json:"answer"`
		TopicID    int    `json:"topic_id"`
		CategoryID int    `json:"category_id"`
		Error      string `json:"error,omitempty"`
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		prompt := url.QueryEscape(string(msg))
		resp, err := http.Get(
			"http://localhost:9090/ask?question=" + prompt +
				"&idCategory=" + idCategory +
				"&topic=" + topic +
				"&username=" + username,
		)

		if err != nil {
			h.manager.SendToUser(userId, "AI error")
			continue
		}

		func() {
			defer resp.Body.Close()

			var aiResp AIResponse
			if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
				h.manager.SendToUser(userId, `{"error":"invalid AI response"}`)
				return
			}

			if aiResp.TopicID != 0 {
				topic = strconv.Itoa(aiResp.TopicID)
			}
			if aiResp.CategoryID != 0 {
				idCategory = strconv.Itoa(aiResp.CategoryID)
			}

			payload, _ := json.Marshal(aiResp)
			h.manager.SendToUser(userId, string(payload))
		}()
	}
}
