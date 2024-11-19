package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// SetupRoutesAuth sets up the routes for the authentication feature.
// It will create an "api" group and add the following routes:
// - POST /api/register: Registers a new user.
// - POST /api/login: Logs in a user and returns a JWT token.
// - POST /api/refresh: Returns a new JWT token for the user.
// - GET /api/userInfo: Returns the user's information, protected by the JWT middleware.
// - PUT /api/userUpdate: Updates the user's information, protected by the JWT middleware.
// - GET /api/users/:id/public: Returns the user's public information, protected by the JWT middleware.
// - GET /api/auth/google: Redirects the user to the Google authentication page.
// - GET /api/auth/google/callback: Handles the callback from Google after the user has been authenticated.
// - GET /api/confirm_email: Confirms the user's email address.
func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
	api.Get("/userInfo", middlewares.JWTMiddleware, controller.UserHandler)
	api.Put("/userUpdate", middlewares.JWTMiddleware, controller.UserUpdate)
	api.Get("/users/:id/public", middlewares.JWTMiddleware, controller.GetPublicUserInfoHandler)
	api.Get("/auth/google", controller.GoogleLogin)
	api.Get("/auth/google/callback", controller.GoogleCallback)
	api.Get("/confirm_email", controller.ConfirmEmailHandler)
}

// SetupRoutesMatches sets up the routes for managing matches.
// It will create an "api/matches" group and add the following routes:
// - GET /api/matches/nearby: Retrieves all nearby matches.
// - GET /api/matches: Retrieves all matches.
// - POST /api/matches: Creates a new match.
// - GET /api/matches/:id: Retrieves a match by ID.
// - PUT /api/matches/:id: Updates a match.
// - DELETE /api/matches/:id: Deletes a match.
// - POST /api/matches/:id/join: Adds a player to a match.
// - WS /api/matches/:id/chat: Establishes a WebSocket connection for a match chat.
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

// SetupRoutesMatchePlayers sets up the routes for managing match players.
// It will create an "api/matchesPlayers" group and add the following routes:
// - GET /api/matchesPlayers/:match_id: Retrieves all match players associated
//   with a given match ID.
// - POST /api/matchesPlayers/: Creates a new match player.
// - PUT /api/matchesPlayers/assignTeam: Assigns a team to a match player.
// - DELETE /api/matchesPlayers/:match_player_id: Deletes a match player.
func SetupRoutesMatchePlayers(app *fiber.App, controller *controllers.MatchPlayersController) {
	api := app.Group("/api/matchesPlayers")
	api.Get("/:match_id", controller.GetMatchPlayersByMatchIDHandler)
	api.Post("/", controller.CreateMatchPlayerHandler)
	api.Put("/assignTeam", controller.AssignTeamToPlayerHandler)
	api.Delete("/:match_player_id", controller.DeleteMatchPlayerHandler)
}

// SetupChatRoutes sets up the routes for managing chat messages.
// It will create an "api/chat" group and add the following routes:
// - POST /api/chat/send: Sends a message to a match chat.
// - GET /api/chat/:match_id: Retrieves all messages associated with a given match ID.
func SetupChatRoutes(app *fiber.App, controller *controllers.ChatController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Post("/chat/send", controller.SendMessage)
	api.Get("/chat/:matchID", controller.GetMessages)
}


// SetupOpenAiRoutes sets up the routes for using OpenAI services.
// It will create an "api" group and add the following routes:
// - GET /api/openai/formation/:match_id: Retrieves a suggested formation
//   for a given match ID, based on the statistics of the players in the match.
func SetupOpenAiRoutes(app *fiber.App, controller *controllers.OpenAiController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Get("/openai/formation/:match_id", controller.GetFormationFromAi)
}
