package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type FriendService struct {
	DB          *gorm.DB
	RedisClient *redis.Client
	AuthService *AuthService
	WebSocket   *websocket.Conn
}

func NewFriendService(db *gorm.DB, redisClient *redis.Client, authService *AuthService, webSocket *websocket.Conn) *FriendService {
	return &FriendService{
		DB:          db,
		RedisClient: redisClient,
		AuthService: authService,
		WebSocket:   webSocket,
	}
}

func (s *FriendService) SendFriendRequest(senderID, receiverID string) error {
	ctx := context.Background()
	requestID := fmt.Sprintf("friend_request:%s:%s", senderID, receiverID)
	friendRequest := models.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	data, err := json.Marshal(friendRequest)
	if err != nil {
		return err
	}

	err = s.RedisClient.Set(ctx, requestID, data, 30*24*time.Hour).Err()
	if err != nil {
		return err
	}

	// Envoyer une notification via WebSocket
	notification := map[string]string{
		"type":       "friend_request",
		"senderId":   senderID,
		"receiverId": receiverID,
	}
	notificationData, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	s.WebSocket.WriteMessage(websocket.TextMessage, notificationData)

	return nil
}

func (s *FriendService) AcceptFriendRequest(senderID, receiverID string) error {
	ctx := context.Background()
	requestID := fmt.Sprintf("friend_request:%s:%s", senderID, receiverID)

	data, err := s.RedisClient.Get(ctx, requestID).Result()
	if err != nil {
		return err
	}

	var friendRequest models.FriendRequest
	err = json.Unmarshal([]byte(data), &friendRequest)
	if err != nil {
		return err
	}

	friendRequest.Status = "accepted"
	updatedData, err := json.Marshal(friendRequest)
	if err != nil {
		return err
	}

	err = s.RedisClient.Set(ctx, requestID, updatedData, 30*24*time.Hour).Err()
	if err != nil {
		return err
	}

	friendship1 := models.Friendship{
		UserID:   senderID,
		FriendID: receiverID,
	}
	friendship2 := models.Friendship{
		UserID:   receiverID,
		FriendID: senderID,
	}

	if err := s.DB.Create(&friendship1).Error; err != nil {
		return err
	}
	if err := s.DB.Create(&friendship2).Error; err != nil {
		return err
	}

	// Envoyer une notification via WebSocket
	notification := map[string]string{
		"type":       "friend_request_accepted",
		"senderId":   senderID,
		"receiverId": receiverID,
	}
	notificationData, err := json.Marshal(notification)
	if err != nil {
		return err
	}
	s.WebSocket.WriteMessage(websocket.TextMessage, notificationData)

	return nil
}

func (s *FriendService) GetFriends(userID string) ([]models.Users, error) {
	var friends []models.Users
	if err := s.DB.Joins("JOIN friendships ON friendships.friend_id = users.id").Where("friendships.user_id = ?", userID).Find(&friends).Error; err != nil {
		return nil, err
	}
	return friends, nil
}

func (s *FriendService) SearchUsersByUsername(username string) ([]models.Users, error) {
	var users []models.Users
	query := "%" + username + "%"
	if err := s.DB.Where("username LIKE ?", query).Find(&users).Error; err != nil {
		return nil, err
	}
	fmt.Printf("Search query: %s, Results: %v\n", query, users)
	return users, nil
}
