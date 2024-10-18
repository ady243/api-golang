package services

import (
    "errors"
    "math/rand"
    "time"

    "github.com/ady243/teamup/helpers"
    middlewares "github.com/ady243/teamup/internal/middleware"
    "github.com/ady243/teamup/internal/models"
    "github.com/oklog/ulid/v2"
    "gorm.io/gorm"
)

// AuthService fournit des services d'authentification
type AuthService struct {
    DB *gorm.DB
}


//ici on crée une nouvelle instance de AuthService 
func NewAuthService(db *gorm.DB) *AuthService {
    return &AuthService{
        DB: db,
    }
}

// Register enregistre un nouvel utilisateur
func (s *AuthService) Register(username, email, password, profilePhoto, favoriteSport, location, bio string, birthDate *time.Time, role models.Role, skillLevel string, pac, sho, pas, dri, def, phy, matchesPlayed, matchesWon, goalsScored, behaviorScore int) (models.Users, error) {
    entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
    newID := ulid.MustNew(ulid.Timestamp(time.Now()), entropy)

    hashedPassword, err := helpers.HashPassword(password)
    if err != nil {
        return models.Users{}, err
    }

    user := models.Users{
        ID:            newID.String(), 
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

// Login authentifie un utilisateur et retourne un token JWT et un refreshToken
func (s *AuthService) Login(email, password string) (string, string, error) {
    var user models.Users
    if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return "", "", err
    }

    if !helpers.CheckPasswordHash(password, user.PasswordHash) {
        return "", "", ErrInvalidCredentials
    }

    userID, err := ulid.Parse(user.ID)
    if err != nil {
        return "", "", err
    }
    accessToken, err := middlewares.GenerateToken(userID)
    if err != nil {
        return "", "", err
    }

    userID, err = ulid.Parse(user.ID)
    if err != nil {
        return "", "", err
    }
    refreshToken, err := middlewares.GenerateRefreshToken(userID)
    if err != nil {
        return "", "", err
    }

    return accessToken, refreshToken, nil
}

func(s *AuthService) GetUserByID(id string) (models.Users, error) {
    var user models.Users
    if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
        return models.Users{}, err
    }

    return user, nil
}

// Refresh génère un nouveau accessToken à partir d'un refreshToken valide
func (s *AuthService) Refresh(refreshToken string) (string, error) {
    claims, err := middlewares.ParseToken(refreshToken)
    if err != nil {
        return "", err
    }

    if claims.ExpiresAt < time.Now().Unix() {
        return "", errors.New("refresh token expired")
    }

    accessToken, err := middlewares.GenerateToken(claims.UserID)
    if err != nil {
        return "", err
    }

    return accessToken, nil
}

// UpdateRefreshToken met à jour le token de rafraîchissement pour un utilisateur spécifique
func (s *AuthService) UpdateRefreshToken(email string, refreshToken string) error {
    var user models.Users
    if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
        return err
    }
    user.RefreshToken = refreshToken
    if err := s.DB.Save(&user).Error; err != nil {
        return err
    }

    return nil
}

// Erreurs spécifiques pour le service
var ErrInvalidCredentials = errors.New("invalid credentials")