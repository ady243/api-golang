package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

// SendFriendRequest sends a friend request from the sender to the receiver
func (s *FriendService) SendFriendRequest(senderID, receiverID string) error {
	ctx := context.Background()

	// Construct Redis key to check if the request already exists
	requestID := fmt.Sprintf("friend_request:%s:%s", senderID, receiverID)

	// Check if the friend request already exists in Redis
	exists, err := s.RedisClient.Exists(ctx, requestID).Result()
	if err != nil {
		return fmt.Errorf("failed to check if friend request exists: %w", err)
	}
	if exists > 0 {
		return fmt.Errorf("friend request already exists")
	}

	// Create the friend request data
	friendRequest := models.FriendRequest{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	// Serialize the friend request
	data, err := json.Marshal(friendRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal friend request: %w", err)
	}

	// Store the friend request in Redis with a TTL of 30 days
	err = s.RedisClient.Set(ctx, requestID, data, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store friend request in Redis: %w", err)
	}

	// If WebSocket is available, send a notification
	if s.WebSocket != nil {
		notification := map[string]string{
			"type":       "friend_request",
			"senderId":   senderID,
			"receiverId": receiverID,
		}

		// Serialize the notification
		notificationData, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("failed to marshal notification: %w", err)
		}

		// Send the notification via WebSocket
		err = s.WebSocket.WriteMessage(websocket.TextMessage, notificationData)
		if err != nil {
			log.Printf("Failed to send WebSocket notification: %v", err)
		}
	} else {
		log.Println("WebSocket is nil. Skipping notification.")
	}

	return nil
}

// AcceptFriendRequest accepts a pending friend request
func (s *FriendService) AcceptFriendRequest(senderID, receiverID string) error {
	ctx := context.Background()

	// Construct Redis key to get the friend request
	requestID := fmt.Sprintf("friend_request:%s:%s", senderID, receiverID)

	// Get the friend request data from Redis
	data, err := s.RedisClient.Get(ctx, requestID).Result()
	if err != nil {
		return fmt.Errorf("failed to get friend request from Redis: %w", err)
	}

	// Unmarshal the friend request data
	var friendRequest models.FriendRequest
	err = json.Unmarshal([]byte(data), &friendRequest)
	if err != nil {
		return fmt.Errorf("failed to unmarshal friend request: %w", err)
	}

	// Update the status to accepted
	friendRequest.Status = "accepted"
	updatedData, err := json.Marshal(friendRequest)
	if err != nil {
		return fmt.Errorf("failed to marshal updated friend request: %w", err)
	}

	// Store the updated friend request in Redis
	err = s.RedisClient.Set(ctx, requestID, updatedData, 30*24*time.Hour).Err()
	if err != nil {
		return fmt.Errorf("failed to store updated friend request in Redis: %w", err)
	}

	// Create friendship entries in the database
	friendship1 := models.Friendship{
		UserID:   senderID,
		FriendID: receiverID,
	}
	friendship2 := models.Friendship{
		UserID:   receiverID,
		FriendID: senderID,
	}

	// Insert friendships into the database
	if err := s.DB.Create(&friendship1).Error; err != nil {
		return fmt.Errorf("failed to create friendship 1: %w", err)
	}
	if err := s.DB.Create(&friendship2).Error; err != nil {
		return fmt.Errorf("failed to create friendship 2: %w", err)
	}

	// Send a notification via WebSocket
	if s.WebSocket != nil {
		notification := map[string]string{
			"type":       "friend_request_accepted",
			"senderId":   senderID,
			"receiverId": receiverID,
		}

		// Serialize the notification
		notificationData, err := json.Marshal(notification)
		if err != nil {
			return fmt.Errorf("failed to marshal notification: %w", err)
		}

		// Send the notification via WebSocket
		err = s.WebSocket.WriteMessage(websocket.TextMessage, notificationData)
		if err != nil {
			log.Printf("Failed to send WebSocket notification: %v", err)
		}
	} else {
		log.Println("WebSocket is nil. Skipping notification.")
	}

	return nil
}

// GetFriends retrieves the list of friends for a user
func (s *FriendService) GetFriends(userID string) ([]models.Users, error) {
	var friends []models.Users
	err := s.DB.Joins("JOIN friendships ON friendships.friend_id = users.id").
		Where("friendships.user_id = ?", userID).
		Find(&friends).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get friends from database: %w", err)
	}
	return friends, nil
}

// SearchUsersByUsername searches for users by their username
func (s *FriendService) SearchUsersByUsername(username string) ([]models.Users, error) {
	var users []models.Users
	query := "%" + username + "%"
	err := s.DB.Where("username LIKE ?", query).Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to search users by username: %w", err)
	}
	fmt.Printf("Search query: %s, Results: %v\n", query, users)
	return users, nil
}

// GetFriendRequests retrieves the pending friend requests for a user
func (s *FriendService) GetFriendRequests(userID string) ([]models.FriendRequest, error) {
	var friendRequests []models.FriendRequest
	err := s.DB.Where("receiver_id = ? AND status = ?", userID, "pending").
		Find(&friendRequests).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get friend requests from database: %w", err)
	}
	return friendRequests, nil
}
