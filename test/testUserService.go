package test

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewAuthService(db *gorm.DB) *services.AuthService {
	return services.NewAuthService(db)
}

func setupTestDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	db.AutoMigrate(&models.Users{})
	return db
}

func TestRegister(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService(db)

	birthDate := time.Now().AddDate(-25, 0, 0)
	user, err := authService.Register("testuser", "test@teamUp.com", "password123", "photo.jpg", "soccer", "NY", "bio", &birthDate, models.Player, "beginner", 80, 70, 60, 50, 40, 30, 10, 5, 3, 90)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "test@teamUp.com", user.Email)

	// Vérifier que l'utilisateur est bien inséré dans la base de données
	var dbUser models.Users
	result := db.Where("email = ?", "test@teamUp.com").First(&dbUser)
	assert.NoError(t, result.Error)
	assert.Equal(t, "testuser", dbUser.Username)
	assert.Equal(t, "test@teamUp.com", dbUser.Email)
}

func TestLogin(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService(db)

	birthDate := time.Now().AddDate(-25, 0, 0)
	_, _ = authService.Register("testuser", "test@teamUp.com", "password123", "photo.jpg", "soccer", "NY", "bio", &birthDate, models.Player, "beginner", 80, 70, 60, 50, 40, 30, 10, 5, 3, 90)

	accessToken, refreshToken, err := authService.Login("test@teamUp.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)

	_, _, err = authService.Login("test@teamUp.com", "wrongpassword")
	assert.Error(t, err)
}

func TestRefresh(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService(db)

	birthDate := time.Now().AddDate(-25, 0, 0)
	_, _ = authService.Register("testuser", "test@teamUp.com", "password123", "photo.jpg", "soccer", "NY", "bio", &birthDate, models.Player, "beginner", 80, 70, 60, 50, 40, 30, 10, 5, 3, 90)

	_, refreshToken, _ := authService.Login("test@teamUp.com", "password123")

	newAccessToken, err := authService.Refresh(refreshToken)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccessToken)
}

func TestUpdateRefreshToken(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService(db)

	birthDate := time.Now().AddDate(-25, 0, 0)
	_, _ = authService.Register("testuser", "test@teamUp.com", "password123", "photo.jpg", "soccer", "NY", "bio", &birthDate, models.Player, "beginner", 80, 70, 60, 50, 40, 30, 10, 5, 3, 90)

	newRefreshToken := ulid.Make().String()
	err := authService.UpdateRefreshToken("test@teamUp.com", newRefreshToken)
	assert.NoError(t, err)

	var user models.Users
	db.Where("email = ?", "test@teamUp.com").First(&user)
	assert.Equal(t, newRefreshToken, user.RefreshToken)
}

func TestUpdateUser(t *testing.T) {
	db := setupTestDB()
	authService := NewAuthService(db)

	birthDate := time.Now().AddDate(-25, 0, 0)
	user, _ := authService.Register("testuser", "test@teamUp.com", "password123", "photo.jpg", "soccer", "NY", "bio", &birthDate, models.Player, "beginner", 80, 70, 60, 50, 40, 30, 10, 5, 3, 90)

	updatedUser, err := authService.UpdateUser(
		user.ID,
		"updateduser",
		"updated@teamUp.com",
		"newpassword",
		"newphoto.jpg",
		"basketball",
		"New York",
		"Updated Bio",
		&birthDate,
		models.Player,
		"Advanced",
		85, 80, 75, 90, 65, 95, 20, 10, 15, 110,
	)

	assert.NoError(t, err)
	assert.Equal(t, "updateduser", updatedUser.Username)
	assert.Equal(t, "updated@teamUp.com", updatedUser.Email)
	assert.Equal(t, "newphoto.jpg", updatedUser.ProfilePhoto)
	assert.Equal(t, "basketball", updatedUser.FavoriteSport)
	assert.Equal(t, "New York", updatedUser.Location)
	assert.Equal(t, "Updated Bio", updatedUser.Bio)
	assert.Equal(t, "Advanced", updatedUser.SkillLevel)
	assert.Equal(t, 85, updatedUser.Pac)
	assert.Equal(t, 80, updatedUser.Sho)
	assert.Equal(t, 75, updatedUser.Pas)
	assert.Equal(t, 90, updatedUser.Dri)
	assert.Equal(t, 65, updatedUser.Def)
	assert.Equal(t, 95, updatedUser.Phy)
	assert.Equal(t, 20, updatedUser.MatchesPlayed)
	assert.Equal(t, 10, updatedUser.MatchesWon)
	assert.Equal(t, 15, updatedUser.GoalsScored)
	assert.Equal(t, 110, updatedUser.BehaviorScore)
}
