package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
	api.Get("/userInfo", middlewares.JWTMiddleware, controller.UserHandler)
	api.Put("/userUpdate", middlewares.JWTMiddleware, controller.UserUpdate)
	api.Get("/auth/google", controller.GoogleLogin)
	api.Get("/auth/google/callback", controller.GoogleCallback)
}

func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")
	api.Get("/nearby", controller.GetNearbyMatchesHandler)
	api.Use(middlewares.JWTMiddleware)
	api.Get("/", controller.GetAllMatchesHandler)
	api.Post("/", controller.CreateMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
	api.Post("/:id/join", controller.AddPlayerToMatchHandler)
	api.Get("/:id/chat", websocket.New(controller.ChatWebSocketHandler))

}

func SetupRoutesMatchePlayers(app *fiber.App, controller *controllers.MatchPlayersController) {
	api := app.Group("/api/matchesPlayers")
	api.Get("/:match_id", controller.GetMatchPlayersByMatchIDHandler)
	api.Post("/", controller.CreateMatchPlayerHandler)
	api.Put("/assignTeam", controller.AssignTeamToPlayerHandler)
	api.Delete("/:match_player_id", controller.DeleteMatchPlayerHandler)
}
func SetupChatRoutes(app *fiber.App, controller *controllers.ChatController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Post("/chat/send", controller.SendMessage)
	api.Get("/chat/:matchID", controller.GetMessages)
}
