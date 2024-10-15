package controllers

import (
	"log"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MatchController struct {
	MatchService *services.MatchService
	DB           *gorm.DB
}

func NewMatchController(matchService *services.MatchService, db *gorm.DB) *MatchController {
	return &MatchController{
		MatchService: matchService,
		DB:           db,
	}
}

func (ctrl *MatchController) GetUserByID(userID string) (*models.Users, error) {
	var user models.Users
	if err := ctrl.DB.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ctrl *MatchController) CreateMatchHandler(c *fiber.Ctx) error {
	var req struct {
		OrganizerID     string `json:"organizer_id" binding:"required"`
		RefereeID       string `json:"referee_id"`
		Description     string `json:"description"`
		MatchDate       string `json:"match_date" binding:"required"`
		MatchTime       string `json:"match_time" binding:"required"`
		Address         string `json:"address" binding:"required"`
		NumberOfPlayers int    `json:"number_of_players" binding:"required"`
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

	log.Printf("Request received: %+v", req)

	organizerID, err := ulid.Parse(req.OrganizerID)
	log.Printf("Organizer ID: %s", organizerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid organizer ID"})
	}

	// Vérifie si l'utilisateur existe avec l'ULID parsé
	user, err := ctrl.GetUserByID(organizerID.String()) // Ou selon la signature de votre fonction
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Organizer not found"})
	}

	organizerID = user.ID

	var refereeID *ulid.ULID
	if req.RefereeID != "" {
		id, err := ulid.Parse(req.RefereeID)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid referee ID"})
		}
		refereeID = &id
	}

	match := &models.Matches{
		OrganizerID:     organizerID,
		RefereeID:       refereeID,
		Description:     &req.Description,
		MatchDate:       matchDate,
		MatchTime:       matchTime,
		Address:         req.Address,
		NumberOfPlayers: req.NumberOfPlayers,
		Status:          "upcoming",
	}

	if err := ctrl.MatchService.CreateMatch(match); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(match)
}

func (ctrl *MatchController) GetMatchByIDHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	matchID, err := ulid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	match, err := ctrl.MatchService.GetMatchByID(matchID)
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

	var req struct {
		Description     string `json:"description"`
		MatchDate       string `json:"match_date"`
		MatchTime       string `json:"match_time"`
		Address         string `json:"address"`
		NumberOfPlayers int    `json:"number_of_players"`
		Status          string `json:"status"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	match, err := ctrl.MatchService.GetMatchByID(matchID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match not found"})
	}

	if req.Description != "" {
		match.Description = &req.Description
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

	if err := ctrl.MatchService.UpdateMatch(match); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(match)
}

func (ctrl *MatchController) DeleteMatchHandler(c *fiber.Ctx) error {
	id := c.Params("id")

	matchID, err := ulid.Parse(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	if err := ctrl.MatchService.DeleteMatch(matchID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
