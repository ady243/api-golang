package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutesAuth sets up the routes for the authentication feature.
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

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/userInfo", controller.UserHandler)
	api.Delete("/deleteMyAccount", controller.DeleteUserHandler)
	api.Get("/users/:id/public", controller.GetPublicUserInfoHandler)
	api.Post("/assignRole/:organizerID/:playerID", controller.AssignRefereeRole)
	api.Post("/UpdateUserStatistics", controller.UpdateUserStatistics)
}

// SetupRoutesMatches sets up the routes for managing matches.
func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")

	// All these routes require JWT auth
	api.Use(middlewares.JWTMiddleware)

	api.Get("/nearby", controller.GetNearbyMatchesHandler)
	api.Get("/", controller.GetAllMatchesHandler)
	api.Post("/", controller.CreateMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
	api.Post("/:id/join", controller.AddPlayerToMatchHandler)
	api.Get("/:id/chat", websocket.New(controller.ChatWebSocketHandler))
	api.Get("/organizer/matches", controller.GetMatchByOrganizerIDHandler)
	api.Get("/analyst/matches", controller.GetMatchByRefereeIDHandler)
	api.Put("/assignAsAnalyst/:match_id/:referee_id", controller.PutRefereeIDHandler)
}

// SetupRoutesMatchePlayers sets up the routes for managing match players.
func SetupRoutesMatchePlayers(app *fiber.App, controller *controllers.MatchPlayersController) {
	api := app.Group("/api/matchesPlayers")

	// Require authentication
	api.Use(middlewares.JWTMiddleware)

	api.Get("/:match_id", controller.GetMatchPlayersByMatchIDHandler)
	api.Post("/", controller.CreateMatchPlayerHandler)
	api.Put("/assignTeam", controller.AssignTeamToPlayerHandler)
	api.Delete("/:match_player_id", controller.DeleteMatchPlayerHandler)
}

// SetupChatRoutes sets up the routes for managing chat messages.
func SetupChatRoutes(app *fiber.App, controller *controllers.ChatController) {
	api := app.Group("/api")

	// Require authentication
	api.Use(middlewares.JWTMiddleware)

	api.Post("/chat/send", controller.SendMessage)
	api.Get("/chat/:matchID", controller.GetMessages)
}

// SetupOpenAiRoutes sets up the routes for using OpenAI services.
func SetupOpenAiRoutes(app *fiber.App, controller *controllers.OpenAiController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)

	api.Get("/openai/formation/:match_id", controller.GetFormationFromAi)
}

// SetupRoutesAnalyst sets up the routes for managing Analyst events.
func SetupRoutesAnalyst(app *fiber.App, controller *controllers.AnalystController) {
	api := app.Group("/api/analyst")

	// Require JWT auth for all Analyst routes
	api.Use(middlewares.JWTMiddleware)

	// Create a new event
	api.Post("/events", controller.CreateEventHandler)

	// Get all events for a match
	api.Get("/match/:match_id/events", controller.GetEventsByMatchHandler)

	// Get all events for a player
	api.Get("/player/:player_id/events", controller.GetEventsByPlayerHandler)

	// Update an existing event
	api.Put("/events/:event_id", controller.UpdateEventHandler)

	// Soft-delete an event
	api.Delete("/events/:event_id", controller.DeleteEventHandler)
}
