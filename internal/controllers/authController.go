package controllers

import (
	"time"

	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
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
		Username      string `json:"username" binding:"required"`
		Email         string `json:"email" binding:"required"`
		Password      string `json:"password" binding:"required"`
		ProfilePhoto  string `json:"profilePhoto"`
		FavoriteSport string `json:"favoriteSport"`
		Bio           string `json:"bio"`
		Location      string `json:"location"`
		BirthDate     string `json:"birthDate"`
		Role          string `json:"role"`
		SkillLevel    string `json:"skillLevel"`
		Pac           int    `json:"pac"`
		Sho           int    `json:"sho"`
		Pas           int    `json:"pas"`
		Dri           int    `json:"dri"`
		Def           int    `json:"def"`
		Phy           int    `json:"phy"`
		MatchesPlayed int    `json:"matchesPlayed"`
		MatchesWon    int    `json:"matchesWon"`
		GoalsScored   int    `json:"goalsScored"`
		BehaviorScore int    `json:"behaviorScore"`
		
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Default role to "Player" if not provided
	if req.Role == "" {
		req.Role = "Player"
	}
	role := models.Role(req.Role)

	// Only parse birthDate if it's provided
	var birthDate time.Time
	if req.BirthDate != "" {
		var err error
		birthDate, err = time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birth date format. Use YYYY-MM-DD"})
		}
	}

	user, err := ctrl.AuthService.Register(
		req.Username, req.Email, req.Password, req.ProfilePhoto, req.FavoriteSport,
		req.Location, req.Bio, &birthDate, role, req.SkillLevel, req.Pac, req.Sho,
		req.Pas, req.Dri, req.Def, req.Phy, req.MatchesPlayed, req.MatchesWon,
		req.GoalsScored, req.BehaviorScore,
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// LoginHandler gère la requête de connexion d'un utilisateur
func (ctrl *AuthController) LoginHandler(c *fiber.Ctx) error {
    var req struct {
        Email    string `json:"email" binding:"required"`
        Password string `json:"password" binding:"required"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    accessToken, refreshToken, err := ctrl.AuthService.Login(req.Email, req.Password)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }

    err = ctrl.AuthService.UpdateRefreshToken(req.Email, refreshToken)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken, "refreshToken": refreshToken})
}

// RefreshHandler gère la demande de rafraîchissement du token d'un utilisateur
func (ctrl *AuthController) RefreshHandler(c *fiber.Ctx) error {
    var req struct {
        RefreshToken string `json:"refreshToken" binding:"required"`
    }

    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    newAccessToken, err := ctrl.AuthService.Refresh(req.RefreshToken)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": newAccessToken})
}