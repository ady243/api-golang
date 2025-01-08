package services

import (
    "encoding/json"
    "log"
    "sync"

    "github.com/gofiber/websocket/v2"
)

type WebSocketService struct {
    connections map[string]map[*websocket.Conn]bool 
    broadcast   chan []byte
    mutex       sync.Mutex
}

func NewWebSocketService() *WebSocketService {
    return &WebSocketService{
        connections: make(map[string]map[*websocket.Conn]bool),
        broadcast:   make(chan []byte),
    }
}


func (s *WebSocketService) HandleWebSocket(c *websocket.Conn, matchID string) {
    log.Println("Handling new WebSocket connection for match:", matchID)
    s.mutex.Lock()
    if s.connections[matchID] == nil {
        s.connections[matchID] = make(map[*websocket.Conn]bool)
    }
    s.connections[matchID][c] = true
    s.mutex.Unlock()

    defer func() {
        s.mutex.Lock()
        delete(s.connections[matchID], c)
        s.mutex.Unlock()
        if err := c.Close(); err != nil {
            log.Println("Error closing WebSocket connection:", err)
        }
        log.Println("WebSocket connection closed for match:", matchID)
    }()

    for {
        _, msg, err := c.ReadMessage()
        if err != nil {
            log.Println("Error reading WebSocket message:", err)
            break
        }
        log.Printf("Received message for match %s: %s", matchID, msg)
        s.broadcast <- msg
    }
}

// BroadcastEventToMatch permet de diffuser un événement structuré (JSON) à tous les clients connectés pour un match spécifique
func (s *WebSocketService) BroadcastEventToMatch(matchID string, event interface{}) {
    message, err := json.Marshal(event)
    if err != nil {
        log.Println("Error marshalling event:", err)
        return
    }

    log.Printf("Broadcasting event to match %s: %s", matchID, string(message))
    s.mutex.Lock()
    defer s.mutex.Unlock()

    for conn := range s.connections[matchID] {
        if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
            log.Println("Error sending WebSocket message:", err)
            conn.Close()
            delete(s.connections[matchID], conn)
        }
    }
}

// StartBroadcast diffuse les messages reçus à tous les clients connectés
func (s *WebSocketService) StartBroadcast() {
    for {
        msg := <-s.broadcast
        log.Printf("Broadcasting message: %s", msg)
        s.mutex.Lock()
        for _, conns := range s.connections {
            for conn := range conns {
                if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
                    log.Println("Error writing WebSocket message:", err)
                    conn.Close()
                    delete(conns, conn)
                }
            }
        }
        s.mutex.Unlock()
    }
}