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

func NewMatchController(matchService *services.MatchService, authService *services.AuthService, db *gorm.DB, chatService *services.ChatService) *MatchController {
	return &MatchController{
		MatchService: matchService,
		AuthService:  authService,
		DB:           db,
		ChatService:  chatService,
	}
}

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

	// Enregistre le match dans la base de données
	if err := ctrl.MatchService.CreateMatch(match, user.ID); err != nil {
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

func (ctrl *MatchController) GetMatchByIDHandler(c *fiber.Ctx) error {
    id := c.Params("id")

    matchID, err := ulid.Parse(id)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
    }

    match, err := ctrl.MatchService.GetMatchByID(matchID.String())
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
    }

    return c.Status(fiber.StatusOK).JSON(match)
}

func (ctrl *MatchController) UpdateMatchHandler(c *fiber.Ctx) error {
    id := c.Params("id")

    matchID, err := ulid.Parse(id)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
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

    // Notifier les clients de la mise à jour du statut du match
    if req.Status != nil {
        if err := ctrl.MatchService.NotifyMatchStatusUpdate(match.ID, string(match.Status)); err != nil {
            log.Println("Erreur lors de la notification de la mise à jour du statut du match:", err)
        }
    }

    return c.Status(fiber.StatusOK).JSON(match)
}

// DeleteMatchHandler deletes a match. It requires the user to be the organizer of the match.
// A soft delete is performed by setting the DeletedAt field to the current time.
// The match is removed from the chat and all players are removed from the match as well.
func (ctrl *MatchController) DeleteMatchHandler(c *fiber.Ctx) error {
    id := c.Params("id")
    log.Println("Match ID:", id)

    // Parse l'ID du match
    matchID, err := ulid.Parse(id)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
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
    if err := ctrl.MatchService.DeleteMatch(matchID.String(), userID); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // Notifier les clients de la suppression du match
    if err := ctrl.MatchService.NotifyMatchStatusUpdate(match.ID, "deleted"); err != nil {
        log.Println("Erreur lors de la notification de la suppression du match:", err)
    }

    return c.SendStatus(fiber.StatusNoContent)
}

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

func (ctrl *MatchController) AddPlayerToMatchHandler(c *fiber.Ctx) error {
	matchID := c.Params("id")
	userID := c.Locals("user_id").(string)

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

// Handler pour obtenir les matchs proches basés sur l'utilisateur connecté
func (ctrl *MatchController) GetNearbyMatchesHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := ctrl.AuthService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user"})
	}

	if user.Latitude == 0 || user.Longitude == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User location not set"})
	}

	matches, err := ctrl.MatchService.FindNearbyMatches(user.Latitude, user.Longitude, 6.0)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(matches)
}

// GetMatchByOrganizerIDHandler gets all matches created by the organizer with the given ID.
// The ID is retrieved from the user_id key in the context.
func (ctrl *MatchController) GetMatchByOrganizerIDHandler(c *fiber.Ctx) error {
	organizerID := c.Locals("user_id").(string)
	matches, err := ctrl.MatchService.GetMatchByOrganizerID(organizerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(matches)
}

// GetMatchByRefereeIDHandler gets all matches where the user with the given ID is referee.
// The ID is retrieved from the user_id key in the context.
func (ctrl *MatchController) GetMatchByRefereeIDHandler(c *fiber.Ctx) error {
	refereeID := c.Locals("user_id").(string)
	matches, err := ctrl.MatchService.GetMatchByRefereeID(refereeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(matches)
}

// PutRefereeIDHandler met à jour le referee_id pour un match donné.
func (ctrl *MatchController) PutRefereeIDHandler(c *fiber.Ctx) error {
	var req struct {
		MatchID   string `json:"match_id" binding:"required"`
		RefereeID string `json:"referee_id" binding:"required"`
	}

	// Parse le corps de la requête
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	// Vérification si l'utilisateur est l'organisateur du match
	if !ctrl.AuthService.IsOrganizer(req.MatchID, c.Locals("user_id").(string)) {
		return c.Status(fiber.StatusForbidden).JSON(map[string]interface{}{"error": "Unauthorized"})
	}

	// Mise à jour du referee_id pour le match spécifié
	if err := ctrl.MatchService.PutRefereeID(req.MatchID, req.RefereeID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Referee ID updated successfully"})
}

func (ctrl *MatchController) MatchStatusWebSocketHandler(c *websocket.Conn) {
    ctx := context.Background()
    pubsub := ctrl.RedisClient.Subscribe(ctx, "match_status_updates")
    defer pubsub.Close()

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
}

func (ctrl *MatchController) UpdateMatchStatusesHandler(c *fiber.Ctx) error {
    if err := ctrl.MatchService.UpdateMatchStatuses(); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update match statuses"})
    }
    return c.SendStatus(fiber.StatusOK)
}