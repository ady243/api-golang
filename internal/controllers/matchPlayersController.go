package controllers

import (
	"fmt"
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

// PlayerResponse structure pour contrôler les données exposées
type PlayerResponse struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Role          string `json:"role"`
	ProfilePhoto  string `json:"profile_photo"`
	FavoriteSport string `json:"favorite_sport"`
	Location      string `json:"location"`
	SkillLevel    string `json:"skill_level"`
	Bio           string `json:"bio"`
	Pac           int    `json:"pac"`
	Sho           int    `json:"sho"`
	Pas           int    `json:"pas"`
	Dri           int    `json:"dri"`
	Def           int    `json:"def"`
	Phy           int    `json:"phy"`
	MatchesPlayed int    `json:"matches_played"`
	MatchesWon    int    `json:"matches_won"`
	GoalsScored   int    `json:"goals_scored"`
	BehaviorScore int    `json:"behavior_score"`
}

// NewMatchPlayersController creates a new instance of MatchPlayersController

// @Summary GetMatchPlayersByMatchIDHandler
// @Description Get all players in a match
// @Tags MatchPlayers
// @Accept json
// @Produce json
// @Param match_id path string true "Match ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/matchesPlayers/{match_id} [get]
func NewMatchPlayersController(matchPlayersService *services.MatchPlayersService, authService *services.AuthService, db *gorm.DB) *MatchPlayersController {
	return &MatchPlayersController{
		MatchPlayersService: matchPlayersService,
		AuthService:         authService,
		DB:                  db,
	}
}

// @Summary CreateMatchPlayerHandler
// @Description Add a player to a match
// @Tags MatchPlayers
// @Accept json
// @Produce json
// @Param match_id body string true "Match ID"
// @Param player_id body string true "Player ID"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/matchesPlayers [post]
func (ctrl *MatchPlayersController) GetMatchPlayersByMatchIDHandler(c *fiber.Ctx) error {
	matchID := c.Params("match_id")

	// Retrieve match players
	matchPlayers, err := ctrl.MatchPlayersService.GetMatchPlayersByMatchID(matchID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Match players not found for match ID: " + matchID})
	}

	// Build player statistics response
	playerStats := make([]PlayerResponse, len(matchPlayers))
	for i, matchPlayer := range matchPlayers {
		playerStats[i] = PlayerResponse{
			ID:            matchPlayer.Player.ID,
			Username:      matchPlayer.Player.Username,
			ProfilePhoto:  matchPlayer.Player.ProfilePhoto,
			FavoriteSport: matchPlayer.Player.FavoriteSport,
			Location:      matchPlayer.Player.Location,
			SkillLevel:    matchPlayer.Player.SkillLevel,
			Bio:           matchPlayer.Player.Bio,
			Pac:           matchPlayer.Player.Pac,
			Sho:           matchPlayer.Player.Sho,
			Pas:           matchPlayer.Player.Pas,
			Dri:           matchPlayer.Player.Dri,
			Def:           matchPlayer.Player.Def,
			Phy:           matchPlayer.Player.Phy,
			MatchesPlayed: matchPlayer.Player.MatchesPlayed,
			MatchesWon:    matchPlayer.Player.MatchesWon,
			GoalsScored:   matchPlayer.Player.GoalsScored,
			BehaviorScore: matchPlayer.Player.BehaviorScore,
		}
	}

	// Return response with player statistics
	return c.JSON(map[string]interface{}{
		"players": playerStats,
	})
}

// CreateMatchPlayerHandler adds a player to a match
func (ctrl *MatchPlayersController) CreateMatchPlayerHandler(c *fiber.Ctx) error {
	var req struct {
		MatchID  string `json:"match_id" binding:"required"`
		PlayerID string `json:"player_id" binding:"required"`
	}

	// Parse the request JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid request format"})
	}

	matchID, err := ulid.Parse(req.MatchID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid match ID"})
	}

	// Validate player ID
	playerID, err := ulid.Parse(req.PlayerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid player ID format"})
	}

	// Verify if player exists
	player, err := ctrl.AuthService.GetUserByID(playerID.String())
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Player not found"})
	}

	// Verify if match exists
	if err := ctrl.MatchPlayersService.CheckMatchExists(matchID.String()); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Match not found"})
	}

	// Generate a new ID for the match player
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	matchPlayersID := ulid.MustNew(ulid.Timestamp(t), entropy).String()

	// Create a new MatchPlayers instance
	matchPlayer := &models.MatchPlayers{
		ID:       matchPlayersID,
		MatchID:  matchID.String(),
		PlayerID: player.ID,
	}

	// Add the player to the match
	if err := ctrl.MatchPlayersService.CreateMatchPlayers(matchPlayer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": "Could not add player to match: " + err.Error()})
	}

	// Prepare the response for the added match player
	matchPlayerResponse := map[string]interface{}{
		"match_player": matchPlayer,
		"player": PlayerResponse{
			ID:            player.ID,
			Username:      player.Username,
			ProfilePhoto:  player.ProfilePhoto,
			FavoriteSport: player.FavoriteSport,
			Location:      player.Location,
			SkillLevel:    player.SkillLevel,
			Bio:           player.Bio,
			Pac:           player.Pac,
			Sho:           player.Sho,
			Pas:           player.Pas,
			Dri:           player.Dri,
			Def:           player.Def,
			Phy:           player.Phy,
			MatchesPlayed: player.MatchesPlayed,
			MatchesWon:    player.MatchesWon,
			GoalsScored:   player.GoalsScored,
			BehaviorScore: player.BehaviorScore,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(matchPlayerResponse)
}

func (ctrl *MatchPlayersController) AssignTeamToPlayerHandler(c *fiber.Ctx) error {
	var req struct {
		MatchPlayerID string `json:"match_player_id" binding:"required"`
		TeamNumber    int    `json:"team_number" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": err.Error()})
	}

	if req.TeamNumber != 1 && req.TeamNumber != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Invalid team number"})
	}

	matchPlayer, err := ctrl.MatchPlayersService.GetMatchPlayerByID(req.MatchPlayerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Match player not found"})
	}

	// Vérification si l'utilisateur est l'organisateur du match
	if !ctrl.AuthService.IsOrganizer(matchPlayer.MatchID, c.Locals("user_id").(string)) {
		return c.Status(fiber.StatusForbidden).JSON(map[string]interface{}{"error": "Unauthorized"})
	}

	// Attribuer l'équipe
	matchPlayer.TeamNumber = &req.TeamNumber

	if err := ctrl.MatchPlayersService.UpdateMatchPlayer(matchPlayer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(matchPlayer)
}

func (ctrl *MatchPlayersController) DeleteMatchPlayerHandler(c *fiber.Ctx) error {
	matchPlayerID := c.Params("match_player_id")

	// Récupérer les informations sur le joueur du match
	matchPlayer, err := ctrl.MatchPlayersService.GetMatchPlayerByID(matchPlayerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Match player not found"})
	}

	// Récupérer l'utilisateur connecté
	userID := c.Locals("user_id").(string)

	// Vérifier si l'utilisateur est soit l'organisateur, soit le joueur lui-même
	if userID != matchPlayer.PlayerID && !ctrl.AuthService.IsOrganizer(matchPlayer.MatchID, userID) {
		return c.Status(fiber.StatusForbidden).JSON(map[string]interface{}{"error": "Unauthorized"})
	}

	// Mettre à jour le champ DeletedAt avec l'heure actuelle
	now := time.Now()
	matchPlayer.DeletedAt = &now

	if err := ctrl.MatchPlayersService.UpdateMatchPlayer(matchPlayer); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{"message": "Player marked as removed from match"})
}

func (ctrl *MatchPlayersController) GetMatchesByPlayerIDHandler(c *fiber.Ctx) error {
	playerID := c.Params("player_id")
	if playerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": "Player ID is required"})
	}

	// Retrieve matches by player ID
	matchPlayers, err := ctrl.MatchPlayersService.GetMatchesByPlayerID(playerID)
	if err != nil {
		fmt.Printf("Error retrieving matches for player %s: %v\n", playerID, err)
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "Player not found in MatchPlayers"})
	}

	// Check if matches were found
	if len(matchPlayers) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(map[string]interface{}{"error": "No matches found for this player"})
	}

	// Retrieve match details for each match player
	var matches []models.Matches
	for _, matchPlayer := range matchPlayers {
		var match models.Matches
		if err := ctrl.DB.Where("id = ?", matchPlayer.MatchID).First(&match).Error; err != nil {
			fmt.Printf("Error retrieving match details for match ID %s: %v\n", matchPlayer.MatchID, err)
			continue
		}
		matches = append(matches, match)
	}

	// Return the matches
	return c.Status(fiber.StatusOK).JSON(map[string]interface{}{
		"match_count": len(matches),
		"data":        matches,
	})
}
