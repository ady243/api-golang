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

func SetupRoutesMatches(app *fiber.App, controller *controllers.MatchController) {
	api := app.Group("/api/matches")
	api.Post("/", controller.CreateMatchHandler)      // Créer un nouveau match
	api.Get("/:id", controller.GetMatchByIDHandler)   // Récupérer un match par ID
	api.Put("/:id", controller.UpdateMatchHandler)    // Mettre à jour un match existant
	api.Delete("/:id", controller.DeleteMatchHandler) // Supprimer un match par ID
}
