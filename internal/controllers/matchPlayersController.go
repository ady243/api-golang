package controllers

import (
	"math/rand"
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"gorm.io/gorm"
)

type MatchPlayersController struct {
	MatchPlayersService *services.MatchPlayersService
	AuthService         *services.AuthService
	DB                  *gorm.DB
}

func NewMatchPlayersController(matchPlayersService *services.MatchPlayersService, authService *services.AuthService, db *gorm.DB) *MatchPlayersController {
	return &MatchPlayersController{
		MatchPlayersService: matchPlayersService,
		AuthService:         authService,
		DB:                  db,
	}
}

func (ctrl *MatchPlayersController) GetMatchPlayersByMatchIDHandler(c *fiber.Ctx) error {
	matchID := c.Params("match_id")
	matchPlayers, err := ctrl.MatchPlayersService.GetMatchPlayersByMatchID(matchID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match players not found"})
	}
	return c.JSON(matchPlayers)
}

func (ctrl *MatchPlayersController) CreateMatchPlayerHandler(c *fiber.Ctx) error {
	var req struct {
		MatchID  string `json:"match_id" binding:"required"`
		PlayerID string `json:"player_id" binding:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	matchID, err := ulid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	playerID, err := ulid.Parse(req.PlayerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid player ID"})
	}

	// Vérifie si l'utilisateur existe avec l'ULID parsé
	player, err := ctrl.AuthService.GetUserByID(playerID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Player not found"})
	}
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0) // Utiliser une source d'entropie
	matchPlayersID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	matchPlayer := &models.MatchPlayers{
		ID:       matchPlayersID,
		MatchID:  matchID.String(),
		PlayerID: player.ID,
	}

	if err := ctrl.MatchPlayersService.CreateMatchPlayers(matchPlayer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(matchPlayer)
}

func (ctrl *MatchPlayersController) AssignTeamToPlayerHandler(c *fiber.Ctx) error {
	var req struct {
		MatchPlayerID string `json:"match_player_id" binding:"required"`
		TeamNumber    int    `json:"team_number" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.TeamNumber != 1 && req.TeamNumber != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid team number"})
	}

	matchPlayer, err := ctrl.MatchPlayersService.GetMatchPlayerByID(req.MatchPlayerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match player not found"})
	}

	// Vérification si l'utilisateur est l'organisateur du match
	if !ctrl.AuthService.IsOrganizer(matchPlayer.MatchID, c.Locals("user_id").(string)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Attribuer l'équipe
	matchPlayer.TeamNumber = &req.TeamNumber

	if err := ctrl.MatchPlayersService.UpdateMatchPlayer(matchPlayer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(matchPlayer)
}
