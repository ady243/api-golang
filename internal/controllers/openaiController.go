package controllers

import (
	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
)

type OpenAiController struct {
    openAi             *services.OpenAIService
    matchPlayersService *services.MatchPlayersService
}

func NewOpenAiController(openAi *services.OpenAIService, matchPlayersService *services.MatchPlayersService) *OpenAiController {
    return &OpenAiController{
        openAi:             openAi,
        matchPlayersService: matchPlayersService,
    }
}

func (ctrl *OpenAiController) GetFormationFromAi(c *fiber.Ctx) error {
    matchID := c.Params("match_id")

    // Récupérer les joueurs du match
    matchPlayers, err := ctrl.matchPlayersService.GetMatchPlayersByMatchID(matchID)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Match players not found for match ID: " + matchID})
    }

    var players []models.Users
    for _, matchPlayer := range matchPlayers {
        players = append(players, matchPlayer.Player)
    }

    // Suggérer une formation basée sur les statistiques des joueurs
    formation, err := ctrl.openAi.SuggestFormations(players)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }
    return c.JSON(fiber.Map{
        "formation": formation,
    })
}