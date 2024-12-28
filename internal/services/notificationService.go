package services

import (
	"context"
	"encoding/json"
	"log"

	"github.com/go-redis/redis/v8"
)

type Notification struct {
	ReceiverId string `json:"receiverId"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

type NotificationService struct {
	redisClient *redis.Client
	broadcast   chan Notification
}

func NewNotificationService(redisClient *redis.Client, broadcast chan Notification) *NotificationService {
	return &NotificationService{
		redisClient: redisClient,
		broadcast:   broadcast,
	}
}

func (s *NotificationService) SendWebSocketNotification(receiverId, title, message string) error {
	notification := Notification{
		ReceiverId: receiverId,
		Title:      title,
		Message:    message,
	}

	// Publish notification to Redis
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	if err := s.redisClient.Publish(context.Background(), "notifications", data).Err(); err != nil {
		return err
	}

	// Send notification via WebSocket
	s.broadcast <- notification
	return nil
}

func (s *NotificationService) ListenForNotifications() {
	pubsub := s.redisClient.Subscribe(context.Background(), "notifications")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var notification Notification
		if err := json.Unmarshal([]byte(msg.Payload), &notification); err != nil {
			log.Printf("Failed to unmarshal notification: %v", err)
			continue
		}
		s.broadcast <- notification
	}
}
