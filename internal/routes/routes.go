package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutesAuth sets up the routes for the authentication feature.
// @Summary Setup authentication routes
// @Description Setup routes for user authentication
// @Tags Auth
func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")

	// Routes that do not require authentication
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
	api.Get("/auth/google", controller.GoogleLogin)
	api.Get("/users", controller.GetUsersHandler)
	api.Get("/auth/google/callback", controller.GoogleCallback)
	api.Get("/confirm_email", controller.ConfirmEmailHandler)
	api.Put("/userUpdate", middlewares.JWTMiddleware, controller.UserUpdate)
	// Routes that require authenticationdeja

	api.Use(middlewares.JWTMiddleware)
	api.Get("/userInfo", controller.UserHandler)
	api.Delete("/deleteMyAccount", controller.DeleteUserHandler)
	api.Get("/users/:id/public", controller.GetPublicUserInfoHandler)
}

// SetupRoutesMatches sets up the routes for managing matches.
// @Summary Setup match routes
// @Description Setup routes for managing matches
// @Tags Matches
func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")

	// Routes that do not require authentication
	api.Get("/nearby", controller.GetNearbyMatchesHandler)

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/", controller.GetAllMatchesHandler)
	api.Post("/", controller.CreateMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
	api.Post("/:id/join", controller.AddPlayerToMatchHandler)
	api.Get("/:id/chat", websocket.New(controller.ChatWebSocketHandler))
}

// SetupRoutesMatchePlayers sets up the routes for managing match players.
// @Summary Setup match players routes
// @Description Setup routes for managing match players
// @Tags MatchPlayers
func SetupRoutesMatchePlayers(app *fiber.App, controller *controllers.MatchPlayersController) {
	api := app.Group("/api/matchesPlayers")

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/:match_id", controller.GetMatchPlayersByMatchIDHandler)
	api.Post("/", controller.CreateMatchPlayerHandler)
	api.Put("/assignTeam", controller.AssignTeamToPlayerHandler)
	api.Delete("/:match_player_id", controller.DeleteMatchPlayerHandler)
}

// SetupChatRoutes sets up the routes for managing chat messages.
// @Summary Setup chat routes
// @Description Setup routes for managing chat messages
// @Tags Chat
func SetupChatRoutes(app *fiber.App, controller *controllers.ChatController) {
	api := app.Group("/api")

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Post("/chat/send", controller.SendMessage)
	api.Get("/chat/:matchID", controller.GetMessages)
}

// SetupOpenAiRoutes sets up the routes for using OpenAI services.
// @Summary Setup OpenAI routes
// @Description Setup routes for using OpenAI services
// @Tags OpenAI
func SetupOpenAiRoutes(app *fiber.App, controller *controllers.OpenAiController) {
	api := app.Group("/api")

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/openai/formation/:match_id", controller.GetFormationFromAi)
}
