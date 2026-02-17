package handler

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		clients: make(map[string]*websocket.Conn),
	}
}

func (m *WebSocketManager) AddClient(userId string, conn *websocket.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.clients[userId] = conn
}

func (m *WebSocketManager) RemoveClient(userId string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.clients, userId)
}

func (m *WebSocketManager) SendToUser(userId string, message string) {
	m.mu.RLock()
	conn, ok := m.clients[userId]
	m.mu.RUnlock()

	if ok {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}

func (m *WebSocketManager) Broadcast(message string) {
	m.mu.RLock()
	for _, conn := range m.clients {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
	m.mu.RUnlock()
}

func (m *WebSocketManager) ClientExists(userId string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.clients[userId]
	return exists
}
