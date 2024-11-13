package models

import (
	"time"
)

type Status string

const (
	Upcoming  Status = "upcoming"
	Ongoing   Status = "ongoing"
	Completed Status = "completed"
)

type Matches struct {
	ID              string     `json:"id" gorm:"primaryKey;type:varchar(26)"`   // ID du match
	OrganizerID     string     `json:"organizer_id" gorm:"not null"`            // Référence vers l'ID de l'organisateur (Users.id)
	Organizer       Users      `json:"organizer" gorm:"foreignKey:OrganizerID"` // Clé étrangère vers Users
	RefereeID       *string    `json:"referee_id" gorm:"null"`                  // ID de l'arbitre, nullable (Users.id)
	Referee         *Users     `json:"referee" gorm:"foreignKey:RefereeID"`     // Clé étrangère vers Users, nullable
	Description     *string    `json:"description" gorm:"null"`                 // Description du match, nullable
	MatchDate       time.Time  `json:"date" gorm:"not null"`                    // Date du match
	MatchTime       time.Time  `json:"time" gorm:"not null"`                    // Heure du match
	Address         string     `json:"address" gorm:"not null"`                 // Adresse du match
	NumberOfPlayers int        `json:"number_of_players" gorm:"not null"`       // Nombre de joueurs
	ScoreTeam1      int        `json:"score_team_1" gorm:"default:0"`           // Score de l'équipe 1
	ScoreTeam2      int        `json:"score_team_2" gorm:"default:0"`           // Score de l'équipe 2
	Status          Status     `json:"status"`                                  // Statut du match
	Latitude        float64    `json:"latitude" gorm:"not null"`                // Latitude du match
	Longitude       float64    `json:"longitude" gorm:"not null"`               // Longitude du match
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`        // Date de création
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`        // Date de mise à jour
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`                 // Date de suppression (soft delete)
}
