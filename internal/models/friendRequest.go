package models

import (
	"time"
)

type FriendRequest struct {
	ID         string    `json:"id" gorm:"primaryKey;type:varchar(26)"`
	SenderId   string    `json:"sender_id" gorm:"type:varchar(26)"`
	ReceiverId string    `json:"receiver_id" gorm:"type:varchar(26)"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`

	Sender   Users `json:"sender" gorm:"foreignKey:SenderId"`
	Receiver Users `json:"receiver" gorm:"foreignKey:ReceiverId"`
}
