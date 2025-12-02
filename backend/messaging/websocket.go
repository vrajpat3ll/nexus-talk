package main

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// ClientConnection represents a connected websocket client
type ClientConnection struct {
	UserID string
	Conn   *websocket.Conn
}

// ConnectionManager handles all active WebSocket connections
type ConnectionManager struct {
	sync.RWMutex
	clients map[string]map[*websocket.Conn]bool // UserID -> WebSocket connections
	threads map[string][]string                 // ThreadID -> UserIDs
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		clients: make(map[string]map[*websocket.Conn]bool),
		threads: make(map[string][]string),
	}
}

func (cm *ConnectionManager) AddClient(userID string, conn *websocket.Conn) {
	cm.Lock()
	defer cm.Unlock()

	if _, exists := cm.clients[userID]; !exists {
		cm.clients[userID] = make(map[*websocket.Conn]bool)
	}
	cm.clients[userID][conn] = true
	log.Printf("New client connected: userID=%s", userID)
}

func (cm *ConnectionManager) RemoveClient(userID string, conn *websocket.Conn) {
	cm.Lock()
	defer cm.Unlock()

	if conns, exists := cm.clients[userID]; exists {
		delete(conns, conn)
		if len(conns) == 0 {
			delete(cm.clients, userID)
		}
	}
	log.Printf("Client disconnected: userID=%s", userID)
}

func (cm *ConnectionManager) SendToUser(userID string, message []byte) {
	cm.RLock()
	defer cm.RUnlock()

	if conns, exists := cm.clients[userID]; exists {
		for conn := range conns {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("Error sending message to user %s: %v", userID, err)
				conn.Close()
			} else {
				log.Printf("Sent message to user %s", userID)
			}
		}
	}
}

func (cm *ConnectionManager) SendToThread(threadID string, message []byte, excludeUserID string) {
	cm.RLock()
	userIDs := cm.threads[threadID]
	cm.RUnlock()

	log.Printf("Broadcasting to thread %s (members: %d)", threadID, len(userIDs))
	for _, userID := range userIDs {
		if userID != excludeUserID {
			cm.SendToUser(userID, message)
		}
	}
}

func (cm *ConnectionManager) UpdateThreadMembers(threadID string, userIDs []string) {
	cm.Lock()
	defer cm.Unlock()
	cm.threads[threadID] = userIDs
}
