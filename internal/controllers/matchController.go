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
	"strconv"
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

func NewMatchController(matchService *services.MatchService, authService *services.AuthService, db *gorm.DB, chatService *services.ChatService) *MatchController {
	return &MatchController{
		MatchService: matchService,
		AuthService:  authService,
		DB:           db,
		ChatService:  chatService,
	}
}

// @Summary GetAllMatches
// @Description Get all matches
// @Tags Matches
// @Accept json
// @Produce json
// @Success 200 {object} []models.Matches
// @Failure 500 {object} fiber.Map {"error": "Could not retrieve matches"}
// @Router /api/matches [get]
func (ctrl *MatchController) GetAllMatchesHandler(c *fiber.Ctx) error {
	// Appelle le service pour récupérer tous les matchs
	matches, err := ctrl.MatchService.GetAllMatches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve matches"})
	}

	// Prépare une liste de matchs avec les informations filtrées de l'organisateur
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

	// Retourne les matchs filtrés en réponse
	return c.Status(fiber.StatusOK).JSON(filteredMatches)
}

// Fonction pour obtenir les coordonnées GPS à partir d'une adresse
// @param address: Adresse à géocoder
// @param apiKey: Clé API Google Maps
// @return lat: Latitude
// @return lng: Longitude
// @return error: Erreur si la géolocalisation a échoué
func getCoordinates(address, apiKey string) (float64, float64, error) {
	// URL encode l'adresse
	address = url.QueryEscape(address)

	// Crée l'URL de requête
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/geocode/json?address=%s&key=%s", address, apiKey)

	// Effectue la requête GET à l'API
	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	// Décode la réponse JSON
	var geoResp GeoResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return 0, 0, err
	}

	// Vérifie si des résultats ont été trouvés
	if len(geoResp.Result) == 0 {
		return 0, 0, fmt.Errorf("no results found for address: %s", address)
	}

	// Récupère les coordonnées GPS
	lat := geoResp.Result[0].Geometry.Location.Lat
	lng := geoResp.Result[0].Geometry.Location.Lng
	return lat, lng, nil
}

// @Tags Matches
// @Accept json
// @Produce json
// @Param organizer_id body string true "Organizer ID"
// @Param referee_id body string false "Referee ID"
// @Param description body string false "Description"
// @Param match_date body string true "Match Date"
// @Param match_time body string true "Match Time"
// @Param address body string true "Address"
// @Param number_of_players body int true "Number of Players"
// @Success 201 {object} models.Matches
// @Failure 400 {object} fiber.Map {"error": "Invalid request"}
// @Failure 500 {object} fiber.Map {"error": "Could not create match"}
// @Router /api/matches [p
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
		Address         string  `json:"address" binding:"required"`
		NumberOfPlayers int     `json:"number_of_players" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Conversion de la date et de l'heure
	matchDate, err := time.Parse("2006-01-02", req.MatchDate)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match date format. Use YYYY-MM-DD"})
	}

	matchTime, err := time.Parse("15:04", req.MatchTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match time format. Use HH:MM"})
	}

	organizerID, err := ulid.Parse(req.OrganizerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid organizer ID"})
	}

	// Récupère l'utilisateur (organisateur) en base de données
	user, err := ctrl.AuthService.GetUserByID(organizerID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Organizer not found"})
	}

	// Génère un nouvel ID de match
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	matchID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// Récupérer la latitude et longitude de l'adresse
	googleApiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	lat, lng, err := getCoordinates(req.Address, googleApiKey)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": fmt.Sprintf("Failed to geocode address: %v", err)})
	}

	// Crée un nouveau match avec les informations fournies
	match := &models.Matches{
		ID:              matchID,
		OrganizerID:     user.ID,
		Description:     req.Description,
		MatchDate:       matchDate,
		MatchTime:       matchTime,
		Address:         req.Address,
		NumberOfPlayers: req.NumberOfPlayers,
		Status:          models.Upcoming,
		Latitude:        lat,
		Longitude:       lng,
	}

	// Gestion de l'arbitre si présent
	if req.RefereeID != nil {
		refereeID, err := ulid.Parse(*req.RefereeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referee ID"})
		}
		refereeIDStr := refereeID.String()
		match.RefereeID = &refereeIDStr
	}

	// Créer le match
	// Enregistre le match dans la base de données
	if err := ctrl.MatchService.CreateMatch(match); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Crée un objet simplifié pour l'organisateur à inclure dans la réponse
	organizer := fiber.Map{
		"id":            user.ID,
		"username":      user.Username,
		"profile_photo": user.ProfilePhoto,
	}

	// Retourne le match avec les infos de l'organisateur
	matchWithOrganizer := fiber.Map{
		"id":                match.ID,
		"organizer_id":      match.OrganizerID,
		"organizer":         organizer,
		"referee_id":        match.RefereeID,
		"description":       match.Description,
		"match_date":        match.MatchDate,
		"match_time":        match.MatchTime,
		"address":           match.Address,
		"number_of_players": match.NumberOfPlayers,
		"status":            match.Status,
		"created_at":        match.CreatedAt,
		"updated_at":        match.UpdatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(matchWithOrganizer)
}

// @Summary GetMatchByID
// @Description Get a match by its ID
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} models.Matches
// @Failure 404 {object} fiber.Map {"error": "Match not found"}
// @Failure 400 {object} fiber.Map {"error": "Invalid match ID"}
// @Router /api/matches/{id} [get]
func (ctrl *MatchController) GetMatchByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	matchID, err := ulid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "A) Invalid match ID"})
	}

	match, err := ctrl.MatchService.GetMatchByID(matchID.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	return c.Status(fiber.StatusOK).JSON(match)
}

// @Summary UpdateMatch
// @Description Update a match
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param referee_id body string false "Referee ID"
// @Param description body string false "Description"
// @Param match_date body string false "Match Date"
// @Param match_time body string false "Match Time"
// @Param address body string false "Address"
// @Param number_of_players body int false "Number of Players"
// @Param status body string false "Status"
// @Success 200 {object} models.Matches
// @Failure 404 {object} fiber.Map {"error": "Match not found"}
// @Failure 400 {object} fiber.Map {"error": "Invalid match ID"}
// @Failure 401 {object} fiber.Map {"error": "You are not authorized to update this match"}
// @Router /api/matches/{id} [put]
func (ctrl *MatchController) UpdateMatchHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	matchID, err := ulid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "B) Invalid match ID"})
	}

	// Récupérer l'ID de l'utilisateur connecté via le middleware JWT
	userID := c.Locals("user_id").(string)

	var req struct {
		RefereeID       *string `json:"referee_id"`
		Description     *string `json:"description"`
		MatchDate       string  `json:"match_date"`
		MatchTime       string  `json:"match_time"`
		Address         string  `json:"address"`
		NumberOfPlayers int     `json:"number_of_players"`
		Status          *string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	match, err := ctrl.MatchService.GetMatchByID(matchID.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	// Vérifier si l'utilisateur connecté est l'organisateur du match
	if match.OrganizerID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to update this match"})
	}

	// Mise à jour des champs du match si les données sont fournies
	if req.RefereeID != nil {
		refereeID, err := ulid.Parse(*req.RefereeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referee ID"})
		}
		refereeIDStr := refereeID.String()
		match.RefereeID = &refereeIDStr
	}

	if req.Description != nil {
		match.Description = req.Description
	}
	if req.MatchDate != "" {
		matchDate, err := time.Parse("2006-01-02", req.MatchDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match date format"})
		}
		match.MatchDate = matchDate
	}
	if req.MatchTime != "" {
		matchTime, err := time.Parse("15:04", req.MatchTime)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match time format"})
		}
		match.MatchTime = matchTime
	}
	if req.Address != "" {
		match.Address = req.Address
	}
	if req.NumberOfPlayers != 0 {
		match.NumberOfPlayers = req.NumberOfPlayers
	}
	if req.Status != nil {
		match.Status = models.Status(*req.Status)
	}

	if err := ctrl.MatchService.UpdateMatch(match); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(match)
}

// @Summary DeleteMatch
// @Description Delete a match
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 204 {object} string ""
// @Failure 404 {object} fiber.Map {"error": "Match not found"}
// @Failure 401 {object} fiber.Map {"error": "You are not authorized to delete this match"}
// @Router /api/matches/{id} [delete]
func (ctrl *MatchController) DeleteMatchHandler(c *fiber.Ctx) error {
	id := c.Params("id")
	log.Println("Match ID:", id)

	// Parse l'ID du match
	matchID, err := ulid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "C) Invalid match ID"})
	}

	// Récupère l'ID de l'utilisateur connecté
	userID := c.Locals("user_id").(string)

	// Récupère le match à partir de son ID
	match, err := ctrl.MatchService.GetMatchByID(matchID.String())
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	// Vérifie si l'utilisateur connecté est l'organisateur du match
	if match.OrganizerID != userID {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "You are not authorized to delete this match"})
	}

	// Effectue la suppression douce via le service
	if err := ctrl.MatchService.DeleteMatch(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// @Summary ChatWebSocketHandler
// @Description Handle WebSocket connections for chat
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Router /api/matches/{id}/chat [get]
func (ctrl *MatchController) ChatWebSocketHandler(c *websocket.Conn) {
	matchID := c.Params("id")
	userID := c.Locals("user_id").(string)

	// Vérifier si l'utilisateur est dans le match
	if err := ctrl.MatchService.IsUserInMatch(matchID, userID); err != nil {
		c.Close()
		log.Println("User not in match:", err)
		return
	}

	room := "match:" + matchID
	ctx := context.Background()
	pubsub := ctrl.RedisClient.Subscribe(ctx, room)
	defer pubsub.Close()

	// Goroutine pour écouter les messages de Redis
	go func() {
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Erreur de réception de message dans Redis : %v", err)
				break
			}
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Printf("Erreur d'envoi de message WebSocket : %v", err)
				break
			}
		}
	}()

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Erreur de lecture de message WebSocket : %v", err)
			break
		}

		// Construire un message avec métadonnées
		fullMessage := fiber.Map{
			"playerId":  userID,
			"message":   string(message),
			"timestamp": time.Now().Format(time.RFC3339),
		}

		// Publier le message au format JSON
		msgJSON, _ := json.Marshal(fullMessage)
		ctrl.RedisClient.Set(ctx, "chat:"+matchID+":"+userID, msgJSON, time.Hour*24*7)
		ctrl.RedisClient.Publish(ctx, room, string(msgJSON))
	}
}

// @Summary AddPlayerToMatch
// @Description Add a player to a match
// @Tags Matches
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} fiber.Map {"status": "Joined match and chat"}
// @Failure 400 {object} fiber.Map {"error": "Could not join match"}
// @Failure 409 {object} fiber.Map {"error": "User already in match"}
// @Failure 500 {object} fiber.Map {"error": "Could not join chat"}
// @Router /api/matches/{id}/join [post]
func (ctrl *MatchController) AddPlayerToMatchHandler(c *fiber.Ctx) error {
	matchID := c.Params("id")
	userID := c.Locals("user_id").(string) // Assurez-vous que l'ID utilisateur est disponible dans le contexte

	// Vérifier si l'utilisateur est déjà dans le match
	if err := ctrl.MatchService.IsUserInMatch(matchID, userID); err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "User already in match"})
	}

	// Ajoute l'utilisateur au match
	if err := ctrl.MatchService.JoinMatch(matchID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Inscrit l'utilisateur au chat du match
	if err := ctrl.ChatService.AddUserToChat(matchID, userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not join chat"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "Joined match and chat"})
}

// @Summary GetNearbyMatches
// @Description Get nearby matches
// @Tags Matches
// @Accept json
// @Produce json
// @Param lat query float64 true "Latitude"
// @Param lon query float64 true "Longitude"
// @Success 200 {object} []models.Matches
// @Failure 400 {object} fiber.Map {"error": "Invalid latitude"}
// @Failure 400 {object} fiber.Map {"error": "Invalid longitude"}
// @Failure 500 {object} fiber.Map {"error": "Could not retrieve nearby matches"}
// @Router /api/matches/nearby [get]
func (ctrl *MatchController) GetNearbyMatchesHandler(c *fiber.Ctx) error {
	lat, err := strconv.ParseFloat(c.Query("lat"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid latitude"})
	}

	lon, err := strconv.ParseFloat(c.Query("lon"), 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid longitude"})
	}

	matches, err := ctrl.MatchService.FindNearbyMatches(lat, lon, 6.0) // La distance de recherche est de 6 km
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(matches)
}
