package models

import "time"

type MatchPlayers struct {
	ID         string     `json:"id" gorm:"primaryKey;type:varchar(26)"`
	MatchID    string     `json:"match_id" gorm:"not null"` // Référence au match
	PlayerID   string     `json:"player_id" gorm:"not null"` // Référence à l'utilisateur
	TeamNumber *int       `json:"team_number" gorm:"null"` // Numéro de l'équipe, nullable
	Position   *string    `json:"position" gorm:"null"`     // Position du joueur, nullable
	CreatedAt  time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  *time.Time `json:"deleted_at" gorm:"index"`

	Player Users `json:"player" gorm:"foreignKey:PlayerID"` // Lien vers le modèle Users
	Match  Matches `json:"match" gorm:"foreignKey:MatchID"`   // Lien vers le modèle Matches
}
