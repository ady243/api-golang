package services

import (
	"errors"
	"time"

	"github.com/ady243/teamup/helpers"
	"github.com/ady243/teamup/internal/middleware"
	"github.com/ady243/teamup/internal/models"
	"gorm.io/gorm"
)

// AuthService fournit des services d'authentification
type AuthService struct{
    DB *gorm.DB
}

//On crée une nouvelle instance de AuthService pour gérer les services d'authentification avec la base de données
func NewAuthService(db *gorm.DB) *AuthService{
    return &AuthService{
        DB: db,
    }
}

// Register enregistre un nouvel utilisateur
func (s *AuthService) Register(username, email, password, profilePhoto, favoriteSport, location, bio string, birthDate time.Time, role models.Role, skillLevel string, pac, sho, pas, dri, def, phy, matchesPlayed, matchesWon, goalsScored, behaviorScore int) (models.Users, error) {
    hashedPassword, err := helpers.HashPassword(password)
    if err != nil {
        return models.Users{}, err
    }

    user := models.Users{
        Username:      username,
        Email:         email,
        Location:      location,
        ProfilePhoto:  profilePhoto,
        Bio:           bio,
        FavoriteSport: favoriteSport,
        PasswordHash:  hashedPassword,
        BirthDate:     birthDate,
        Role:          role,
        SkillLevel:    skillLevel,
        Pac:           pac,
        Sho:           sho,
        Pas:           pas,
        Dri:           dri,
        Def:           def,
        Phy:           phy,
        MatchesPlayed: matchesPlayed,
        MatchesWon:    matchesWon,
        GoalsScored:   goalsScored,
        BehaviorScore: behaviorScore,
    }

    if err := s.DB.Create(&user).Error; err != nil {
        return models.Users{}, err
    }

    return user, nil
}

// Login authentifie un utilisateur et retourne un token JWT
func (s *AuthService) Login(email, password string) (string, error) {
    var user models.Users
    if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return "", err
    }

    if !helpers.CheckPasswordHash(password, user.PasswordHash) {
        return "", ErrInvalidCredentials
    }

    token, err := middlewares.GenerateToken(user.ID)
    if err != nil {
        return "", err
    }

    return token, nil
}

// Erreurs spécifiques pour le service
var ErrInvalidCredentials = errors.New("invalid credentials")
