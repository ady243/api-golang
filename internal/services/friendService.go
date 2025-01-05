package services

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

type FriendService struct {
	DB               *gorm.DB
	AuthService      *AuthService
	WebSocketService *WebSocketService
}

func NewFriendService(db *gorm.DB, authService *AuthService, webSocketService *WebSocketService) *FriendService {
	return &FriendService{
		DB:               db,
		AuthService:      authService,
		WebSocketService: webSocketService,
	}
}

// SendFriendRequest sends a friend request from the sender to the receiver
func (s *FriendService) SendFriendRequest(senderId, receiverId string) error {
	log.Printf("Sending friend request from %s to %s", senderId, receiverId)

	// Check if the friend request already exists in the database
	var existingRequest models.FriendRequest
	if err := s.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderId, receiverId, receiverId, senderId).First(&existingRequest).Error; err == nil {
		log.Printf("Friend request already exists in database for %s to %s", senderId, receiverId)
		return fmt.Errorf("friend request already exists")
	}

	// Create the friend request data
	friendRequest := models.FriendRequest{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	// Store the friend request in the database
	if err := s.DB.Create(&friendRequest).Error; err != nil {
		log.Printf("Failed to create friend request in database: %v", err)
		return fmt.Errorf("failed to create friend request in database: %w", err)
	}
	log.Printf("Stored friend request in database for %s to %s", senderId, receiverId)

	// Send notification via WebSocket
	notification := map[string]string{
		"type":       "friend_request",
		"senderId":   senderId,
		"receiverId": receiverId,
	}

	// Serialize the notification
	notificationData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send the notification via WebSocketService
	s.WebSocketService.broadcast <- notificationData

	return nil
}

// AcceptFriendRequest accepts a pending friend request
func (s *FriendService) AcceptFriendRequest(senderId, receiverId string) error {
	log.Printf("Accepting friend request from %s to %s", senderId, receiverId)

	// Check if the friend request exists in the database
	var friendRequest models.FriendRequest
	if err := s.DB.Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).First(&friendRequest).Error; err != nil {
		log.Printf("Friend request not found for %s to %s: %v", senderId, receiverId, err)
		return fmt.Errorf("friend request not found: %w", err)
	}

	log.Printf("Friend request found: %v", friendRequest)

	// Update the status of the friend request
	friendRequest.Status = "accepted"
	if err := s.DB.Save(&friendRequest).Error; err != nil {
		log.Printf("Failed to update friend request status: %v", err)
		return fmt.Errorf("failed to update friend request status: %w", err)
	}

	log.Printf("Accepted friend request from %s to %s", senderId, receiverId)

	// Send notification via WebSocket
	notification := map[string]string{
		"type":       "friend_request_accepted",
		"senderId":   senderId,
		"receiverId": receiverId,
	}

	// Serialize the notification
	notificationData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send the notification via WebSocketService
	s.WebSocketService.broadcast <- notificationData

	return nil
}

// DeclineFriendRequest declines a pending friend request
func (s *FriendService) DeclineFriendRequest(senderId, receiverId string) error {
	log.Printf("Declining friend request from %s to %s", senderId, receiverId)

	// Check if the friend request exists in the database
	var friendRequest models.FriendRequest
	if err := s.DB.Where("sender_id = ? AND receiver_id = ?", senderId, receiverId).First(&friendRequest).Error; err != nil {
		log.Printf("Friend request not found for %s to %s: %v", senderId, receiverId, err)
		return fmt.Errorf("friend request not found: %w", err)
	}

	log.Printf("Friend request found: %v", friendRequest)

	// Update the status of the friend request
	friendRequest.Status = "declined"
	if err := s.DB.Save(&friendRequest).Error; err != nil {
		log.Printf("Failed to update friend request status: %v", err)
		return fmt.Errorf("failed to update friend request status: %w", err)
	}

	log.Printf("Declined friend request from %s to %s", senderId, receiverId)

	// Send notification via WebSocket
	notification := map[string]string{
		"type":       "friend_request_declined",
		"senderId":   senderId,
		"receiverId": receiverId,
	}

	// Serialize the notification
	notificationData, err := json.Marshal(notification)
	if err != nil {
		log.Printf("Failed to marshal notification: %v", err)
		return fmt.Errorf("failed to marshal notification: %w", err)
	}

	// Send the notification via WebSocketService
	s.WebSocketService.broadcast <- notificationData

	return nil
}

// GetFriendRequests retrieves the pending friend requests for a user
func (s *FriendService) GetFriendRequests(userID string) ([]models.FriendRequest, error) {
	log.Printf("Retrieving friend requests for user %s", userID)
	var friendRequests []models.FriendRequest
	err := s.DB.Preload("Sender").Preload("Receiver").Where("receiver_id = ? AND status = ?", userID, "pending").Find(&friendRequests).Error
	if err != nil {
		log.Printf("Failed to get friend requests from database: %v", err)
		return nil, fmt.Errorf("failed to get friend requests from database: %w", err)
	}
	log.Printf("Retrieved %d friend requests for user %s", len(friendRequests), userID)
	return friendRequests, nil
}

// GetFriends retrieves the list of friends for a user
func (s *FriendService) GetFriends(userID string) ([]models.Users, error) {
	log.Printf("Retrieving friends for user %s", userID)
	var friends []models.Users

	// Retrieve friends where the user is the sender
	err := s.DB.Joins("JOIN friend_requests ON friend_requests.receiver_id = users.id").
		Where("friend_requests.sender_id = ? AND friend_requests.status = ?", userID, "accepted").
		Find(&friends).Error
	if err != nil {
		log.Printf("Failed to get friends from database where user is sender: %v", err)
		return nil, fmt.Errorf("failed to get friends from database where user is sender: %w", err)
	}

	// Retrieve friends where the user is the receiver
	var friendsAsReceiver []models.Users
	err = s.DB.Joins("JOIN friend_requests ON friend_requests.sender_id = users.id").
		Where("friend_requests.receiver_id = ? AND friend_requests.status = ?", userID, "accepted").
		Find(&friendsAsReceiver).Error
	if err != nil {
		log.Printf("Failed to get friends from database where user is receiver: %v", err)
		return nil, fmt.Errorf("failed to get friends from database where user is receiver: %w", err)
	}

	// Combine both lists
	friends = append(friends, friendsAsReceiver...)

	log.Printf("Retrieved %d friends for user %s", len(friends), userID)
	return friends, nil
}

// SearchUsersByUsername searches for users by their username
func (s *FriendService) SearchUsersByUsername(username string) ([]models.Users, error) {
	log.Printf("Searching users by username: %s", username)
	var users []models.Users
	query := "%" + username + "%"
	err := s.DB.Where("username LIKE ?", query).Find(&users).Error
	if err != nil {
		log.Printf("Failed to search users by username: %v", err)
		return nil, fmt.Errorf("failed to search users by username: %w", err)
	}
	log.Printf("Search query: %s, Results: %v", query, users)
	return users, nil
}

func (s *FriendService) AreFriends(userID1, userID2 string) (bool, error) {
	var count int64
	err := s.DB.Model(&models.FriendRequest{}).
		Where("(sender_id = ? AND receiver_id = ? AND status = ?) OR (sender_id = ? AND receiver_id = ? AND status = ?)",
			userID1, userID2, "accepted", userID2, userID1, "accepted").
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check friendship status: %w", err)
	}
	return count > 0, nil
}
