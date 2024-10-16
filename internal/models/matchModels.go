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
	Organizer       Users      `json:"organizer" gorm:"foreignKey:OrganizerID"` // Définition de la clé étrangère vers Users
	RefereeID       *string    `json:"referee_id" gorm:"null"`                  // Référence vers l'ID de l'arbitre, nullable (Users.id)
	Referee         *Users     `json:"referee" gorm:"foreignKey:RefereeID"`     // Clé étrangère vers Users, nullable
	Description     *string    `json:"description" gorm:"null"`                 // Description du match, nullable
	MatchDate       time.Time  `json:"date" gorm:"not null"`                    // Date du match
	MatchTime       time.Time  `json:"time" gorm:"not null"`                    // Heure du match
	Address         string     `json:"address" gorm:"not null"`                 // Adresse du match
	NumberOfPlayers int        `json:"number_of_players" gorm:"not null"`       // Nombre de joueurs
	ScoreTeam1      int        `json:"score_team_1" gorm:"default:0"`           // Score de l'équipe 1, par défaut 0
	ScoreTeam2      int        `json:"score_team_2" gorm:"default:0"`           // Score de l'équipe 2, par défaut 0
	Status          Status     `json:"status"`                                  // Statut du match
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`        // Date de création du match
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`        // Date de mise à jour
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`                 // Date de suppression (soft delete)
}
