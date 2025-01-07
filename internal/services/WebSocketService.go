package services

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gofiber/websocket/v2"
)

type WebSocketService struct {
	connections map[*websocket.Conn]bool
	broadcast   chan []byte
	mutex       sync.Mutex
}

// NewWebSocketService initialise un nouveau service WebSocket
func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		connections: make(map[*websocket.Conn]bool),
		broadcast:   make(chan []byte),
	}
}

// HandleWebSocket gère une nouvelle connexion WebSocket
func (s *WebSocketService) HandleWebSocket(c *websocket.Conn) {
	log.Println("Handling new WebSocket connection")
	s.mutex.Lock()
	s.connections[c] = true
	s.mutex.Unlock()

	defer func() {
		s.mutex.Lock()
		delete(s.connections, c)
		s.mutex.Unlock()
		if err := c.Close(); err != nil {
			log.Println("Error closing WebSocket connection:", err)
		}
		log.Println("WebSocket connection closed")
	}()

	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading WebSocket message:", err)
			break
		}
		log.Printf("Received message: %s", msg)
		s.broadcast <- msg
	}
}

// StartBroadcast diffuse les messages reçus à tous les clients connectés
func (s *WebSocketService) StartBroadcast() {
	for {
		msg := <-s.broadcast
		log.Printf("Broadcasting message: %s", msg)
		s.mutex.Lock()
		for conn := range s.connections {
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Error writing WebSocket message:", err)
				conn.Close()
				delete(s.connections, conn)
			}
		}
		s.mutex.Unlock()
	}
}

// BroadcastEvent permet de diffuser un événement structuré (JSON) à tous les clients connectés
func (s *WebSocketService) BroadcastEvent(event interface{}) {
	message, err := json.Marshal(event)
	if err != nil {
		log.Println("Error marshalling event:", err)
		return
	}

	log.Printf("Broadcasting event: %s", string(message))
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for conn := range s.connections {
		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error sending WebSocket message:", err)
			conn.Close()
			delete(s.connections, conn)
		}
	}
}
