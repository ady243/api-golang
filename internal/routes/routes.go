package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	"github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
}

func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")
	api.Post("/", controller.CreateMatchHandler)
	api.Get("/:id", controller.GetMatchByIDHandler)
	api.Use(middlewares.JWTMiddleware2)
	api.Put("/:id", controller.UpdateMatchHandler)
	api.Delete("/:id", controller.DeleteMatchHandler)
}
