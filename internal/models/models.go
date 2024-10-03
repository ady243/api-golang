package models

import (
	"time"

	"github.com/oklog/ulid/v2"
)

type Role string

const (
	Player        Role = "player"
	Referee       Role = "referee"
	Administrator Role = "administrator"
)

type Users struct {
	ID           ulid.ULID  `json:"id" gorm:"primaryKey"`
	Username     string     `json:"username"`
	Email        string     `json:"email" gorm:"unique"`
	PasswordHash string     `json:"password_hash"`
	UpdatedAt    time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt    *time.Time `json:"deleted_at" gorm:"index"`

	BirthDate     *time.Time `json:"birth_date"`
	Role          Role       `json:"role"`
	ProfilePhoto  string     `json:"profile_photo"`
	FavoriteSport string     `json:"favorite_sport"`
	Location      string     `json:"location"`
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
}

type Matches struct {
	ID              ulid.ULID  `json:"id" gorm:"primaryKey"`                                                                  // ID du match
	OrganizerID     ulid.ULID  `json:"organizer_id" gorm:"not null"`                                                          // Référence vers l'ID de l'organisateur (Users.id)
	Organizer       Users      `json:"organizer" gorm:"foreignKey:OrganizerID"`                                               // Définition de la clé étrangère vers Users
	RefereeID       *ulid.ULID `json:"referee_id" gorm:"null"`                                                                // Référence vers l'ID de l'arbitre, nullable (Users.id)
	Referee         *Users     `json:"referee" gorm:"foreignKey:RefereeID"`                                                   // Clé étrangère vers Users, nullable
	Description     *string    `json:"description" gorm:"null"`                                                               // Description du match, nullable
	MatchDate       time.Time  `json:"date" gorm:"not null"`                                                                  // Date du match
	MatchTime       time.Time  `json:"time" gorm:"not null"`                                                                  // Heure du match
	Address         string     `json:"address" gorm:"not null"`                                                               // Adresse du match
	NumberOfPlayers int        `json:"number_of_players" gorm:"not null"`                                                     // Nombre de joueurs
	ScoreTeam1      int        `json:"score_team_1" gorm:"default:0"`                                                         // Score de l'équipe 1, par défaut 0
	ScoreTeam2      int        `json:"score_team_2" gorm:"default:0"`                                                         // Score de l'équipe 2, par défaut 0
	Status          string     `json:"status" gorm:"type:enum('upcoming','ongoing','completed');default:'upcoming';not null"` // Statut du match
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`                                                      // Date de création du match
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`                                                      // Date de mise à jour
	DeletedAt       *time.Time `json:"deleted_at" gorm:"index"`                                                               // Date de suppression (soft delete)
}
