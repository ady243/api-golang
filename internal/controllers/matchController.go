package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/joho/godotenv"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MatchController struct {
	MatchService *services.MatchService
	ChatService  *services.ChatService
	RedisClient  *redis.Client
	AuthService  *services.AuthService
	DB           *gorm.DB
}

type GeoResponse struct {
	Result []struct {
		Geometry struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
}

func NewMatchController(matchService *services.MatchService, authService *services.AuthService, db *gorm.DB, chatService *services.ChatService, redisClient *redis.Client) *MatchController {
	return &MatchController{
		MatchService: matchService,
		AuthService:  authService,
		DB:           db,
		ChatService:  chatService,
		RedisClient:  redisClient,
	}
}

func (ctrl *MatchController) GetAllMatchesHandler(c *fiber.Ctx) error {
	matches, err := ctrl.MatchService.GetAllMatches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve matches"})
	}

	var filteredMatches []fiber.Map
	for _, match := range matches {
		organizer := fiber.Map{
			"id":            match.Organizer.ID,
			"username":      match.Organizer.Username,
			"profile_photo": match.Organizer.ProfilePhoto,
		}
		filteredMatches = append(filteredMatches, fiber.Map{
			"id":                match.ID,
			"organizer":         organizer,
			"referee_id":        match.RefereeID,
			"description":       match.Description,
			"date":              match.MatchDate,
			"time":              match.MatchTime,
			"address":           match.Address,
			"number_of_players": match.NumberOfPlayers,
			"status":            match.Status,
			"created_at":        match.CreatedAt,
			"updated_at":        match.UpdatedAt,
		})
	}

	return c.Status(fiber.StatusOK).JSON(filteredMatches)
}

func getCoordinates(address, apiKey string) (float64, float64, error) {
	address = url.QueryEscape(address)
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	var geoResp GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, err
	}

	if len(geoResp.Result) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	lat := geoResp.Result[0].Geometry.Location.Lat
	lng := geoResp.Result[0].Geometry.Location.Lng
	return lat, lng, nil
}

func (ctrl *MatchController) CreateMatchHandler(c *fiber.Ctx) error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	var req struct {
		OrganizerID     string  `json:"organizer_id" binding:"required"`
		RefereeID       *string `json:"referee_id"`
		Description     *string `json:"description"`
		MatchDate       string  `json:"match_date" binding:"required"`
		MatchTime       string  `json:"match_time" binding:"required"`
		EndTime         string  `json:"end_time" binding:"required"`
		Address         string  `json:"address" binding:"required"`
		NumberOfPlayers int     `json:"number_of_players" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	matchDate, err := time.Parse("2006-01-02", req.MatchDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match date format. Use YYYY-MM-DD"})
	}

	matchTime, err := time.Parse("15:04:05", req.MatchTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match time format. Use HH:MM:SS"})
	}

	endTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid end time format. Use HH:MM:SS"})
	}

	organizerID, err := ulid.Parse(req.OrganizerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid organizer ID"})
	}

	user, err := ctrl.AuthService.GetUserByID(organizerID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Organizer not found"})
	}

	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	matchID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	googleApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	lat, lng, err := getCoordinates(req.Address, googleApiKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to geocode address: %v", err)})
	}

	match := &models.Matches{
		ID:              matchID,
		OrganizerID:     user.ID,
		Description:     req.Description,
		MatchDate:       matchDate,
		MatchTime:       matchTime,
		EndTime:         endTime,
		Address:         req.Address,
		NumberOfPlayers: req.NumberOfPlayers,
		Status:          models.Upcoming,
		Latitude:        lat,
		Longitude:       lng,
	}

	if req.RefereeID != nil {
		refereeID, err := ulid.Parse(*req.RefereeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referee ID"})
		}
		refereeIDStr := refereeID.String()
		match.RefereeID = &refereeIDStr
	}

	if err := ctrl.MatchService.CreateMatch(match, user.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	organizer := fiber.Map{
		"id":            user.ID,
		"username":      user.Username,
		"profile_photo": user.ProfilePhoto,
	}

	matchWithOrganizer := fiber.Map{
		"id":                match.ID,
		"organizer_id":      match.OrganizerID,
		"organizer":         organizer,
		"referee_id":        match.RefereeID,
		"description":       match.Description,
		"match_date":        match.MatchDate,
		"match_time":        match.MatchTime,
		"end_time":          match.EndTime,
		"address":           match.Address,
		"number_of_players": match.NumberOfPlayers,
		"status":            match.Status,
		"created_at":        match.CreatedAt,
		"updated_at":        match.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(matchWithOrganizer)
}
