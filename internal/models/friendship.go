package models

type Friendship struct {
	ID       string `json:"id" gorm:"primaryKey;type:varchar(26)"`
	UserID   string `json:"user_id" gorm:"type:varchar(26)"`
	FriendID string `json:"friend_id" gorm:"type:varchar(26)"`
}
