package controllers

import (
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
)

type AnalystController struct {
	AnalystService *services.AnalystService
	AuthService    *services.AuthService // si besoin de vérifier l'utilisateur
	DB             *gorm.DB
}

// NewAnalystController retourne un nouveau contrôleur
func NewAnalystController(analystService *services.AnalystService, authService *services.AuthService, db *gorm.DB) *AnalystController {
	return &AnalystController{
		AnalystService: analystService,
		AuthService:    authService,
		DB:             db,
	}
}

// CreateEventHandler crée un nouvel événement (ex: but, carton, etc.)
func (ctrl *AnalystController) CreateEventHandler(c *fiber.Ctx) error {
	var req struct {
		MatchID   string `json:"match_id" binding:"required"`
		AnalystID string `json:"analyst_id" binding:"required"`
		PlayerID  string `json:"player_id" binding:"required"`
		EventType string `json:"event_type" binding:"required"`
		Minute    int    `json:"minute"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request format"})
	}

	// Vérifier si les ULID sont valides (si vous utilisez ULID pour vos IDs)
	matchID, err := ulid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID format"})
	}
	analystID, err := ulid.Parse(req.AnalystID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid analyst ID format"})
	}
	playerID, err := ulid.Parse(req.PlayerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid player ID format"})
	}

	// Générer un nouvel ID pour l'événement
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	eventID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	event := models.Analyst{
		ID:        eventID,
		MatchID:   matchID.String(),
		AnalystID: analystID.String(),
		PlayerID:  playerID.String(),
		EventType: req.EventType,
		Minute:    req.Minute,
	}

	// Sauvegarde en base de l'événement
	if err := ctrl.AnalystService.CreateEvent(&event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not create event: " + err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(event)
}

// GetEventsByMatchHandler renvoie tous les événements d'un match
func (ctrl *AnalystController) GetEventsByMatchHandler(c *fiber.Ctx) error {
	matchID := c.Params("match_id")

	// Vérifier l'existence du match (optionnel, selon votre logique)
	events, err := ctrl.AnalystService.GetEventsByMatchID(matchID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "No events found for match ID: " + matchID})
	}

	return c.JSON(fiber.Map{"events": events})
}

// GetEventsByPlayerHandler renvoie tous les événements pour un joueur donné
func (ctrl *AnalystController) GetEventsByPlayerHandler(c *fiber.Ctx) error {
	playerID := c.Params("player_id")

	events, err := ctrl.AnalystService.GetEventsByPlayerID(playerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"events": events})
}

// UpdateEventHandler met à jour un événement (par ex. changer le type ou la minute)
func (ctrl *AnalystController) UpdateEventHandler(c *fiber.Ctx) error {
	eventID := c.Params("event_id")

	var req struct {
		EventType *string `json:"event_type"`
		Minute    *int    `json:"minute"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Récupérer l'event existant
	event, err := ctrl.AnalystService.GetEventByID(eventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	// Appliquer les modifications
	if req.EventType != nil {
		event.EventType = *req.EventType
	}
	if req.Minute != nil {
		event.Minute = *req.Minute
	}

	// Sauvegarder
	if err := ctrl.AnalystService.UpdateEvent(event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(event)
}

// DeleteEventHandler supprime un événement (soft delete => on remplit DeletedAt)
func (ctrl *AnalystController) DeleteEventHandler(c *fiber.Ctx) error {
	eventID := c.Params("event_id")

	event, err := ctrl.AnalystService.GetEventByID(eventID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Event not found"})
	}

	now := time.Now()
	event.DeletedAt = &now

	if err := ctrl.AnalystService.UpdateEvent(event); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Event marked as deleted"})
}
