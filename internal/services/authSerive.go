package services

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"github.com/ady243/teamup/helpers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/ady243/teamup/internal/models"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

// AuthService fournit des services d'authentification
type AuthService struct {
	DB                *gorm.DB
	GoogleOauthConfig *oauth2.Config
	ImageService      *ImageService
}

// NewAuthService crée une nouvelle instance de AuthService
func NewAuthService(db *gorm.DB, imageService *ImageService) *AuthService {
	googleOauthConfig := &oauth2.Config{
		ClientID:    os.Getenv("GOOGLE_CLIENT_ID"),
		RedirectURL: os.Getenv("GOOGLE_REDIRECT_URI"),
		Scopes:      []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:    google.Endpoint,
	}

	return &AuthService{
		DB:                db,
		GoogleOauthConfig: googleOauthConfig,
		ImageService:      imageService,
	}
}

// GetUserByEmail recherche un utilisateur par email
func (s *AuthService) GetUserByEmail(email string) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return models.Users{}, err
	}
	return user, nil
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
	accessToken, err := middlewares.GenerateToken(userID, user.Role)
    if err != nil {
        return "", "", err
    }

	userID, err = ulid.Parse(user.ID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := middlewares.GenerateRefreshToken(userID, user.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *AuthService) GetUserByID(id string) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

// UpdateUser met à jour les informations d'un utilisateur
func (s *AuthService) UpdateUser(id, username, email, password, profilePhoto, favoriteSport, location, bio string, birthDate *time.Time, role models.Role, skillLevel string, pac, sho, pas, dri, def, phy, matchesPlayed, matchesWon, goalsScored, behaviorScore int) (models.Users, error) {
	var user models.Users
	if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
		return models.Users{}, err
	}

	if username != "" {
		user.Username = username
	}
	if email != "" {
		user.Email = email
	}
	if password != "" {
		hashedPassword, err := helpers.HashPassword(password)
		if err != nil {
			return models.Users{}, err
		}
		user.PasswordHash = hashedPassword
	}
	if profilePhoto != "" {
		user.ProfilePhoto = profilePhoto
	}
	if favoriteSport != "" {
		user.FavoriteSport = favoriteSport
	}
	if location != "" {
		user.Location = location
	}
	if bio != "" {
		user.Bio = bio
	}
	if birthDate != nil {
		user.BirthDate = birthDate
	}
	if role != "" {
		user.Role = role
	}
	if skillLevel != "" {
		user.SkillLevel = skillLevel
	}
	if pac != 0 {
		user.Pac = pac
	}
	if sho != 0 {
		user.Sho = sho
	}
	if pas != 0 {
		user.Pas = pas
	}
	if dri != 0 {
		user.Dri = dri
	}
	if def != 0 {
		user.Def = def
	}
	if phy != 0 {
		user.Phy = phy
	}
	if matchesPlayed != 0 {
		user.MatchesPlayed = matchesPlayed
	}
	if matchesWon != 0 {
		user.MatchesWon = matchesWon
	}
	if goalsScored != 0 {
		user.GoalsScored = goalsScored
	}
	if behaviorScore != 0 {
		user.BehaviorScore = behaviorScore
	}

	if err := s.DB.Save(&user).Error; err != nil {
		return models.Users{}, err
	}

	return user, nil
}

func (s *AuthService) GetPublicUserInfo(id string) (map[string]interface{}, error) {
    var user models.Users
    if err := s.DB.Where("id = ?", id).First(&user).Error; err != nil {
        return nil, err
    }

    publicInfo := map[string]interface{}{
        "username":       user.Username,
        "matchesPlayed":  user.MatchesPlayed,
        "matchesWon":     user.MatchesWon,
        "goalsScored":    user.GoalsScored,
        "behaviorScore":  user.BehaviorScore,
		"profilePhoto":   user.ProfilePhoto,
		"favoriteSport":  user.FavoriteSport,
		"location":       user.Location,
		"bio":            user.Bio,
		"skillLevel":     user.SkillLevel,
		"sho":            user.Sho,
		"pas":            user.Pas,
		"dri":            user.Dri,
		"def":            user.Def,
		"phy":            user.Phy,

    }

    return publicInfo, nil
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

	accessToken, err := middlewares.GenerateToken(claims.UserID, claims.Role)
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

// Vérifie si le user est bien l'organisateur du match
func (s *AuthService) IsOrganizer(matchID, userID string) bool {
	var match models.Matches
	if err := s.DB.Where("id = ? AND organizer_id = ?", matchID, userID).First(&match).Error; err != nil {
		return false
	}
	return true
}
