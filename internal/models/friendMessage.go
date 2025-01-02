package models

import (
	"time"
)

type Message struct {
	ID         uint      `gorm:"primaryKey"`
	SenderID   string    `gorm:"not null"`
	ReceiverID string    `gorm:"not null"`
	Content    string    `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
}
