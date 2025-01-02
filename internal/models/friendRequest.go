package models

import "time"

type FriendRequest struct {
	ID         uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	SenderId   string     `json:"sender_id" gorm:"type:varchar(26)"`
	ReceiverId string     `json:"receiver_id" gorm:"type:varchar(26)"`
	Status     string     `json:"status" gorm:"type:varchar(10)"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`

	Sender   Users `gorm:"foreignKey:SenderId" json:"sender"`
	Receiver Users `gorm:"foreignKey:ReceiverId" json:"receiver"`
}
