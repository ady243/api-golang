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
	api.Get("/auth/google/callback", controller.GoogleCallback)
	api.Get("/confirm_email", controller.ConfirmEmailHandler)
	api.Get("/users", controller.GetUsersHandler)

	api.Put("/userUpdate", middlewares.JWTMiddleware, controller.UserUpdate)

	// Routes that require authentication
	api.Use(middlewares.JWTMiddleware)
	api.Get("/userInfo", controller.UserHandler)
	api.Put("/userUpdate", controller.UserUpdate)
	api.Delete("/deleteMyAccount", controller.DeleteUserHandler)
	api.Get("/users/:id/public", controller.GetPublicUserInfoHandler)
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
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
	api.Post("/:id/join", controller.AddPlayerToMatchHandler)
	api.Post("/:id/leave", controller.LeaveMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Get("/:id/chat", websocket.New(controller.ChatWebSocketHandler))
	api.Get("/organizer/matches", controller.GetMatchByOrganizerIDHandler)
	api.Get("/referee/matches", controller.GetMatchByRefereeIDHandler)
	api.Put("/assignAsAnalyst/:match_id/:referee_id", controller.PutRefereeIDHandler)
	api.Get("/status/updates", websocket.New(controller.MatchStatusWebSocketHandler))
	api.Get("/matches/status/updates", websocket.New(controller.MatchStatusWebSocketHandler))
	api.Post("/assign-referee", controller.AssignRefereeHandler)
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
	api.Use(middlewares.JWTMiddleware)

	api.Post("/events", controller.CreateEventHandler)
	api.Get("/match/:match_id/events", controller.GetEventsByMatchHandler)
	api.Get("/player/:player_id/events", controller.GetEventsByPlayerHandler)
	api.Put("/events/:event_id", controller.UpdateEventHandler)
	api.Delete("/events/:event_id", controller.DeleteEventHandler)
}

// SetupRoutesFriend sets up the routes for managing friend requests.
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

// SetupRoutesFriendMessage sets up the routes for managing messages between friends.
func SetupRoutesFriendMessage(app *fiber.App, friendChatController *controllers.FriendChatController) {
	api := app.Group("/api")
	api.Use(middlewares.JWTMiddleware)
	api.Post("/message/send", friendChatController.SendMessage)
	api.Get("/message/messages/:senderID/:receiverID", friendChatController.GetMessages)
}

// SetupRoutesWebSocket sets up the routes for WebSocket connections.
func SetupRoutesWebSocket(app *fiber.App, controller *controllers.WebSocketController) {
	api := app.Group("/api")
	api.Get("/matches/status/updates", websocket.New(controller.WebSocketHandler))
}

func SetupNotificationRoutes(app *fiber.App, notificationController *controllers.NotificationController) {
	api := app.Group("/api")
	api.Post("/send-notification", notificationController.SendPushNotification)
	api.Get("/notifications/:token", notificationController.GetUnreadNotifications)
	api.Post("/notifications/:token/read", notificationController.MarkNotificationsAsRead)
	api.Post("/send-notification", notificationController.SendPushNotification)
}
