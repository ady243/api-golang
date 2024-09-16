package controllers

import (
    "github.com/ady243/teamup/services"
	"github.com/gofiber/fiber/v2"
)

// AuthController fournit les gestionnaires pour les opérations d'authentification
type AuthController struct {
    AuthService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		AuthService: authService,
	}
}

// RegisterHandler gère la requête d'enregistrement d'un nouvel utilisateur
func (ctrl *AuthController) RegisterHandler(c *fiber.Ctx) error {
    var req struct {
        Username string `json:"username" binding:"required"`
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
        ProfilePhoto string `json:"profilePhoto"`
        FavoriteSport string `json:"favoriteSport"`
        Bio string `json:"bio"`
        Location string `json:"location"` 
    }
    
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    user, err := ctrl.AuthService.Register(req.Username, req.Email, req.Password, req.ProfilePhoto, req.FavoriteSport, req.Location, req.Bio)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(user)
}

// LoginHandler gère la requête de connexion d'un utilisateur
func (ctrl *AuthController) LoginHandler(c *fiber.Ctx) error {
    var req struct {
        Email string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    token, err := ctrl.AuthService.Login(req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}