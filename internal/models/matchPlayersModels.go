package models

import "time"

type MatchPlayers struct {
	ID         string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	MatchID    string     `json:"match_id" gorm:"not null"`
	PlayerID   string     `json:"player_id" gorm:"not null"`
	TeamNumber *int       `json:"team_number" gorm:"null"`
	Position   *string    `json:"position" gorm:"null"`
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`
}
