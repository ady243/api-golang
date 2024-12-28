package models

import "time"

type Friendship struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    string    `gorm:"not null"`
	FriendID  string    `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
}

type FriendRequest struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   string    `gorm:"not null"`
	ReceiverID string    `gorm:"not null"`
	Status     string    `gorm:"not null;default:'pending'"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
