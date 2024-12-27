package models

import (
	"time"
)

type Analyst struct {
	ID        string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	MatchID   string     `json:"match_id" gorm:"not null"`   // Référence au match
	AnalystID string     `json:"analyst_id" gorm:"not null"` // Référence à l'utilisateur qui enregistre l'événement
	PlayerID  string     `json:"player_id" gorm:"not null"`  // Référence à l'utilisateur (joueur) concerné
	EventType string     `json:"event_type" gorm:"null"`     // Type de l'événement (ex: goal, yellow_card, etc.)
	Minute    int        `json:"minute"`                     // Minute de l'événement dans le match
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`

	AnalystUser Users   `json:"analyst_user" gorm:"foreignKey:AnalystID"`
	Player      Users   `json:"player"       gorm:"foreignKey:PlayerID"`
	Match       Matches `json:"match"        gorm:"foreignKey:MatchID"`
}
