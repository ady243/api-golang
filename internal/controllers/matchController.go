package controllers

import (
	"log"
	"math/rand"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MatchController struct {
	MatchService *services.MatchService
	AuthService  *services.AuthService
	DB           *gorm.DB
}

func NewMatchController(matchService *services.MatchService, authService *services.AuthService, db *gorm.DB) *MatchController {
	return &MatchController{
		MatchService: matchService,
		AuthService:  authService,
		DB:           db,
	}
}

func (ctrl *MatchController) GetAllMatchesHandler(c *fiber.Ctx) error {
	// Appelle le service pour récupérer tous les matchs
	matches, err := ctrl.MatchService.GetAllMatches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve matches"})
	}

	// Retourne les matchs en réponse
	return c.Status(fiber.StatusOK).JSON(matches)
}

func (ctrl *MatchController) CreateMatchHandler(c *fiber.Ctx) error {
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
    }

    if req.RefereeID != nil {
        refereeID, err := ulid.Parse(*req.RefereeID)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referee ID"})
        }
        refereeIDStr := refereeID.String()
        match.RefereeID = &refereeIDStr
    }

    // Enregistre le match dans la base de données
    if err := ctrl.MatchService.CreateMatch(match); err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // Crée un objet simplifié pour l'organisateur à inclure dans la réponse
    organizer := fiber.Map{
        "id":        user.ID,
        "username":  user.Username,
        "profile_photo": user.ProfilePhoto, 
    }

    // Retourne le match avec les infos de l'organisateur
    matchWithOrganizer := fiber.Map{
        "id":              match.ID,
        "organizer_id":    match.OrganizerID,
        "organizer":       organizer, 
        "referee_id":      match.RefereeID,
        "description":     match.Description,
        "match_date":      match.MatchDate,
        "match_time":      match.MatchTime,
        "address":         match.Address,
        "number_of_players": match.NumberOfPlayers,
        "status":          match.Status,
        "created_at":      match.CreatedAt,
        "updated_at":      match.UpdatedAt,
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
	if err := ctrl.MatchService.DeleteMatch(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent) 
}
