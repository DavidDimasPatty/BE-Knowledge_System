package handler

import (
	"bufio"
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
	// idCategory := c.Query("idCategory")
	// topic := c.Query("topic")
	username := c.Query("username")
	role := c.Query("roleName")
	// isFirst := c.Query("isFirst")

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
		Answer   string `json:"answer"`
		Topic    int    `json:"topic_id"`
		Category int    `json:"category_id"`
		Error    string `json:"error,omitempty"`
	}

	type ClientMessage struct {
		Text       string `json:"text"`
		IsFirst    bool   `json:"isFirst"`
		IdCategory int    `json:"idCategory"`
		Topic      int    `json:"topic"`
	}

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			return
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			h.manager.SendToUser(userId, `{"error":"invalid message format"}`)
			continue
		}

		prompt := url.QueryEscape(clientMsg.Text)
		isFirst := clientMsg.IsFirst
		idCategory := clientMsg.IdCategory
		topic := clientMsg.Topic
		//username := clientMsg.Username
		//prompt := url.QueryEscape(string(msg))
		resp, err := http.Get(
			"http://localhost:9090/ask?question=" + prompt +
				"&idCategory=" + strconv.Itoa(idCategory) +
				"&topic=" + strconv.Itoa(topic) + "&role=" + role +
				"&username=" + username + "&isFirst=" + strconv.FormatBool(isFirst),
		)

		if err != nil {
			h.manager.SendToUser(userId, "Internal error")
			continue
		}

		// func() {
		// 	defer resp.Body.Close()

		// 	var aiResp AIResponse
		// 	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		// 		h.manager.SendToUser(userId, `{"error":"invalid AI response"}`)
		// 		return
		// 	}

		// 	// if aiResp.Topic != "" {
		// 	// 	topic = aiResp.Topic
		// 	// }
		// 	// if aiResp.Category != "" {
		// 	// 	idCategory = aiResp.Category
		// 	// }

		// 	payload, err := json.Marshal(aiResp)
		// 	if err != nil {
		// 		h.manager.SendToUser(userId, `{"error":"encode failed"}`)
		// 		return
		// 	}

		// 	h.manager.SendToUser(userId, string(payload))
		// }()
		func() {
			defer resp.Body.Close()

			reader := bufio.NewReader(resp.Body)

			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					// stream selesai
					break
				}

				// validasi JSON
				var chunk map[string]interface{}
				if err := json.Unmarshal(line, &chunk); err != nil {
					continue
				}

				// forward ke client (React)
				h.manager.SendToUser(userId, string(line))
			}
		}()
	}
}
