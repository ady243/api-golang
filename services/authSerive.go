package  services

import (
    "errors" 
 "github.com/ady243/teamup/helpers"
    "github.com/ady243/teamup/middleware"
    "github.com/ady243/teamup/models"
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
func (s *AuthService) Register(username, email, password, ProfilePhoto, FavoriteSport,location, Bio string) (models.Users, error) {
    hashedPassword, err := helpers.HashPassword(password)
    if err != nil {
        return models.Users{}, err
    }

    user := models.Users{
        Username:     username,
		Email:        email,
		Location:   location,
		ProfilePhoto: ProfilePhoto,
		Bio:  Bio,
		FavoriteSport:  FavoriteSport,
        PasswordHash: hashedPassword,
    }

    if err := s.DB.Create(&user).Error; err != nil {
        return models.Users{}, err
    }

    return user, nil
}
// Remove the unused function addInformation
// func (s *AuthService) addInformation(ProfilePhoto, location, FavoriteSport, Bio string) (models.Users, error) {
// 	
// 	user := models.Users{
// 		ProfilePhoto:  ProfilePhoto,
// 		Location:      location,
// 		FavoriteSport: FavoriteSport,
// 		Bio:           Bio,
// 	}
// 
// 	if err := s.DB.Create(&user).Error; err != nil {
// 		return models.Users{}, err
// 	}
// 
// 	return user, nil
// }

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
