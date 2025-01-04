package services

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/ady243/teamup/internal/models"
    "gorm.io/gorm"
)

const (
    fcmURL = "https://fcm.googleapis.com/fcm/send"
)

type FCMMessage struct {
    To           string                 `json:"to"`
    Notification map[string]string      `json:"notification"`
    Data         map[string]interface{} `json:"data"`
}

type NotificationService struct {
    DB *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
    return &NotificationService{
        DB: db,
    }
}

func sendFCMNotification(serverKey, deviceToken, title, body string, data map[string]interface{}) error {
    fcmMessage := FCMMessage{
        To: deviceToken,
        Notification: map[string]string{
            "title": title,
            "body":  body,
        },
        Data: data,
    }

    jsonData, err := json.Marshal(fcmMessage)
    if err != nil {
        return fmt.Errorf("failed to marshal FCM message: %v", err)
    }

    req, err := http.NewRequest("POST", fcmURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return fmt.Errorf("failed to create HTTP request: %v", err)
    }

    req.Header.Set("Authorization", "key="+serverKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("failed to send HTTP request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("failed to send FCM notification, status code: %d", resp.StatusCode)
    }

    return nil
}

func (s *NotificationService) SendMessageNotification(senderID, receiverID, content string) error {
    var sender, receiver models.Users
    if err := s.DB.First(&sender, "id = ?", senderID).Error; err != nil {
        log.Printf("Failed to get sender: %v", err)
        return fmt.Errorf("failed to get sender: %w", err)
    }

    if err := s.DB.First(&receiver, "id = ?", receiverID).Error; err != nil {
        log.Printf("Failed to get receiver: %v", err)
        return fmt.Errorf("failed to get receiver: %w", err)
    }

    serverKey := os.Getenv("FCM_SERVER_KEY")
    title := "Nouveau message"
    body := fmt.Sprintf("Vous avez reçu un nouveau message de %s", sender.Username)
    data := map[string]interface{}{
        "type":        "message",
        "sender_id":   sender.ID,
        "username":    sender.Username,
        "receiver_id": receiver.ID,
        "message":     content,
    }

    if err := sendFCMNotification(serverKey, receiver.FCMToken, title, body, data); err != nil {
        log.Printf("Failed to send FCM notification: %v", err)
        return fmt.Errorf("failed to send FCM notification: %w", err)
    }

    return nil
}

func (s *NotificationService) SendFriendRequestNotification(senderID, receiverID string) error {
    var sender, receiver models.Users
    if err := s.DB.First(&sender, "id = ?", senderID).Error; err != nil {
        log.Printf("Failed to get sender: %v", err)
        return fmt.Errorf("failed to get sender: %w", err)
    }

    if err := s.DB.First(&receiver, "id = ?", receiverID).Error; err != nil {
        log.Printf("Failed to get receiver: %v", err)
        return fmt.Errorf("failed to get receiver: %w", err)
    }

    serverKey := os.Getenv("FCM_SERVER_KEY")
    title := "Nouvelle demande d'ami"
    body := fmt.Sprintf("Vous avez reçu une demande d'ami de %s", sender.Username)
    data := map[string]interface{}{
        "type":        "friend_request",
        "sender_id":   sender.ID,
        "username":    sender.Username,
        "receiver_id": receiver.ID,
    }

    if err := sendFCMNotification(serverKey, receiver.FCMToken, title, body, data); err != nil {
        log.Printf("Failed to send FCM notification: %v", err)
        return fmt.Errorf("failed to send FCM notification: %w", err)
    }

    return nil
}