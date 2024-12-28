package models

import (
	"time"
)

type Role string

type Users struct {
	ID           string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	Username     string     `json:"username"`
	Email        string     `json:"email" gorm:"unique"`
	PasswordHash string     `json:"password_hash"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	BirthDate     *time.Time `json:"birth_date"`
	ProfilePhoto  string     `json:"profile_photo"`
	FavoriteSport string     `json:"favorite_sport"`
	Location      string     `json:"location"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	SkillLevel    string     `json:"skill_level"`
	Bio           string     `json:"bio"`

	Pac           int `json:"pac"`
	Sho           int `json:"sho"`
	Pas           int `json:"pas"`
	Dri           int `json:"dri"`
	Def           int `json:"def"`
	Phy           int `json:"phy"`
	MatchesPlayed int `json:"matches_played"`
	MatchesWon    int `json:"matches_won"`
	GoalsScored   int `json:"goals_scored"`
	BehaviorScore int `json:"behavior_score"`

	RefreshToken      string     `json:"refresh_token"`
	IsConfirmed       bool       `json:"is_confirmed" gorm:"default:false"`
	ConfirmationToken string     `json:"confirmation_token" gorm:"size:255"`
	TokenExpiresAt    *time.Time `json:"token_expires_at"`

	SentFriendRequests     []FriendRequest `json:"sent_friend_requests" gorm:"foreignKey:SenderId"`
	ReceivedFriendRequests []FriendRequest `json:"received_friend_requests" gorm:"foreignKey:ReceiverId"`
}
