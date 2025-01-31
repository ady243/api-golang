package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

type FriendChatService struct {
	DB                  *gorm.DB
	WebSocketService    *WebSocketService
	NotificationService *NotificationService
}

func NewFriendChatService(db *gorm.DB, webSocketService *WebSocketService, notificationService *NotificationService) *FriendChatService {
	return &FriendChatService{
		DB:                  db,
		WebSocketService:    webSocketService,
		NotificationService: notificationService,
	}
}

func (s *FriendChatService) SendMessage(senderID, receiverID, content string) error {
	log.Printf("Sending message from %s to %s: %s", senderID, receiverID, content)

	// Create the message data
	message := models.Message{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		CreatedAt:  time.Now(),
	}

	// Store the message in the database
	if err := s.DB.Create(&message).Error; err != nil {
		log.Printf("Failed to create message in database: %v", err)
		return fmt.Errorf("failed to create message in database: %w", err)
	}
	log.Printf("Stored message in database for %s to %s", senderID, receiverID)

	// Récupérer le token FCM du destinataire
	var receiver models.Users
	if err := s.DB.Where("id = ?", receiverID).First(&receiver).Error; err != nil {
		log.Printf("Failed to fetch receiver: %v", err)
		return fmt.Errorf("failed to fetch receiver: %w", err)
	}

	// Envoi de notification via Firebase
	err := s.NotificationService.SendPushNotification(
		receiver.FCMToken,
		"Vous avez un nouveau message",
		content,
	)
	if err != nil {
		log.Printf("Failed to send push notification: %v", err)
		return fmt.Errorf("failed to send push notification: %w", err)
	}

	notification := map[string]string{
		"type":       "new_message",
		"senderID":   senderID,
		"receiverID": receiverID,
		"content":    content,
	}

	notificationData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	s.WebSocketService.broadcast <- notificationData

	return nil
}

func (s *FriendChatService) GetMessages(senderID, receiverID string) ([]models.Message, error) {
	log.Printf("Retrieving messages between %s and %s", senderID, receiverID)
	var messages []models.Message
	err := s.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID).Order("created_at asc").Find(&messages).Error
	if err != nil {
		log.Printf("Failed to get messages from database: %v", err)
		return nil, fmt.Errorf("failed to get messages from database: %w", err)
	}
	log.Printf("Retrieved %d messages between %s and %s", len(messages), senderID, receiverID)
	return messages, nil
}
