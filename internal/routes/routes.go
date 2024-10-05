package routes


import (
	"github.com/ady243/teamup/internal/controllers"
		"github.com/gofiber/fiber/v2"
)


func SetupRoutesAuth(app *fiber.App, controller *controllers.AuthController) {
	api := app.Group("/api")
	api.Post("/register", controller.RegisterHandler)
	api.Post("/login", controller.LoginHandler)
	api.Post("/refresh", controller.RefreshHandler)
}