package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ady243/teamup/internal/models"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type Notification struct {
	ReceiverId string `json:"receiverId"`
	Title      string `json:"title"`
	Message    string `json:"message"`
}

type NotificationService struct {
	db               *gorm.DB
	redisClient      *redis.Client
	broadcast        chan Notification
	webSocketService *WebSocketService
}

func NewNotificationService(db *gorm.DB, redisClient *redis.Client, broadcast chan Notification, webSocketService *WebSocketService) *NotificationService {
	return &NotificationService{
		db:               db,
		redisClient:      redisClient,
		broadcast:        broadcast,
		webSocketService: webSocketService,
	}
}

const (
	fcmServerKey = "YOUR_SERVER_KEY" // Remplacez par votre clé serveur FCM
	fcmURL       = "https://fcm.googleapis.com/fcm/send"
)

type FCMMessage struct {
	To           string            `json:"to"`
	Notification FCMNotification   `json:"notification"`
	Data         map[string]string `json:"data"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (ns *NotificationService) sendPushNotification(token, title, body string) error {
	message := FCMMessage{
		To: token,
		Notification: FCMNotification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"title": title,
			"body":  body,
		},
	}

	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonMessage))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "key="+fcmServerKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification, status code: %d", resp.StatusCode)
	}

	return nil
}

func (ns *NotificationService) SendWebSocketNotification(userID, title, message string) error {
	notification := map[string]string{
		"userID":  userID,
		"title":   title,
		"message": message,
	}

	notificationData, err := json.Marshal(notification)
	if err != nil {
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	ns.webSocketService.broadcast <- notificationData

	// Envoyer une notification push via FCM
	token, err := ns.getUserFCMToken(userID)
	if err != nil {
		return fmt.Errorf("failed to get FCM token: %w", err)
	}

	if err := ns.sendPushNotification(token, title, message); err != nil {
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	return nil
}

func (ns *NotificationService) getUserFCMToken(userID string) (string, error) {
	var user models.Users
	if err := ns.db.First(&user, "id = ?", userID).Error; err != nil {
		return "", fmt.Errorf("failed to find user: %w", err)
	}
	if user.FCMToken == "" {
		return "", fmt.Errorf("user does not have a FCM token")
	}
	return user.FCMToken, nil
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
