package routes

import (
	"github.com/ady243/teamup/internal/controllers"
	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/gofiber/fiber/v2"
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
