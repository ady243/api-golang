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
	api.Get("/auth/google/callback", controller.GoogleCallback)
	api.Get("/confirm_email", controller.ConfirmEmailHandler)
	api.Put("/userUpdate", middlewares.JWTMiddleware, controller.UserUpdate)

	// Routes that require authentication

	api.Use(middlewares.JWTMiddleware)
	api.Get("/userInfo", controller.UserHandler)
	api.Get("/users", controller.GetUsersHandler)
	api.Delete("/deleteMyAccount", controller.DeleteUserHandler)
	api.Get("/users/:id/public", controller.GetPublicUserInfoHandler)
	api.Post("/UpdateUserStatistics", controller.UpdateUserStatistics)

}

// SetupRoutesMatches sets up the routes for managing matches.
// @Summary Setup match routes
// @Description Setup routes for managing matches
// @Tags Matches
func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")
	api.Use(middlewares.JWTMiddleware)
	api.Get("/nearby", controller.GetNearbyMatchesHandler)
	api.Get("/", controller.GetAllMatchesHandler)
	api.Post("/", controller.CreateMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
	api.Post("/:id/join", controller.AddPlayerToMatchHandler)
	api.Post("/:id/leave", controller.LeaveMatchHandler)
	api.Get("/:id/chat", websocket.New(controller.ChatWebSocketHandler))
	api.Get("/organizer/matches", controller.GetMatchByOrganizerIDHandler)
	api.Get("/referee/matches", controller.GetMatchByRefereeIDHandler)
	api.Get("/matches/status/updates", websocket.New(controller.MatchStatusWebSocketHandler))
	api.Post("/assign-referee", controller.AssignRefereeHandler)
}

// SetupRoutesMatchePlayers sets up the routes for managing match players.
// It will create an "api/matchesPlayers" group and add the following routes:
//   - GET /api/matchesPlayers/:match_id: Retrieves all match players associated
//     with a given match ID.
//   - POST /api/matchesPlayers/: Creates a new match player.
//   - PUT /api/matchesPlayers/assignTeam: Assigns a team to a match player.
//   - DELETE /api/matchesPlayers/:match_player_id: Deletes a match player.
//   - GET /api/matchesPlayers/player/player_id: get match player by player.
func SetupRoutesMatchePlayers(app *fiber.App, controller *controllers.MatchPlayersController) {
	api := app.Group("/api/matchesPlayers")

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/:match_id", controller.GetMatchPlayersByMatchIDHandler)
	api.Post("/", controller.CreateMatchPlayerHandler)
	api.Put("/assignTeam", controller.AssignTeamToPlayerHandler)
	api.Delete("/:match_player_id", controller.DeleteMatchPlayerHandler)
	api.Get("/player/:player_id", controller.GetMatchesByPlayerIDHandler)
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
// It will create an "api" group and add the following routes:
//   - GET /api/openai/formation/:match_id: Retrieves a suggested formation
//     for a given match ID, based on the statistics of the players in the match.
func SetupOpenAiRoutes(app *fiber.App, controller *controllers.OpenAiController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Get("/openai/formation/:match_id", controller.GetFormationFromAi)
}

// SetupRoutesWebSocket sets up the routes for managing WebSocket connections.
// It will create an "api" group and add the following routes:
//   - GET /api/matches/status/updates: WebSocket route for receiving match status updates.
func SetupRoutesWebSocket(app *fiber.App, controller *controllers.WebSocketController) {
	api := app.Group("/api")
	api.Get("/matches/status/updatess", websocket.New(controller.WebSocketHandler))
}

// SetupRoutesFriend sets up the routes for managing friend requests.
// It will create an "api" group and add the following routes:
//   - POST /api/friend/send: Sends a friend request.
//   - POST /api/friend/accept: Accepts a friend request.
func SetupFriendRoutes(app *fiber.App, friendController *controllers.FriendController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Post("/friend/send", friendController.SendFriendRequest)
	api.Post("/friend/accept", friendController.AcceptFriendRequest)
	api.Post("/friend/decline", friendController.DeclineFriendRequest)
	api.Get("/friend/requests/:userID", friendController.GetFriendRequests)
	api.Get("/friend/:userID", friendController.GetFriends)
	api.Get("/friend/search", friendController.SearchUsersByUsername)
}

func SetupRoutesFriendMessage(app *fiber.App, friendChatController *controllers.FriendChatController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Post("/message/send", friendChatController.SendMessage)
	api.Get("/message/messages/:senderID/:receiverID", friendChatController.GetMessages)
}

func SetupNotificationRoutes(app *fiber.App, notificationController *controllers.NotificationController) {
	app.Post("/send-notification", notificationController.SendPushNotification)
}
