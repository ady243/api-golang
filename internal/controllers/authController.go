package controllers

import (
	"context"
	"encoding/json"
	"math/rand"
	"time"

	middlewares "github.com/ady243/teamup/internal/middleware"
	"github.com/ady243/teamup/internal/models"
	"github.com/ady243/teamup/internal/services"
	"github.com/gofiber/fiber/v2"
	"github.com/oklog/ulid/v2"
	"golang.org/x/oauth2"
)

type AuthController struct {
	AuthService  *services.AuthService
	ImageService *services.ImageService
	MatchService *services.MatchService
}

// NewAuthController creates a new instance of AuthController.
// It requires an AuthService and an ImageService to handle authentication
// and image-related operations, respectively.
func NewAuthController(authService *services.AuthService, imageService *services.ImageService, matchService *services.MatchService) *AuthController {
	return &AuthController{
		AuthService:  authService,
		ImageService: imageService,
		MatchService: matchService,
	}
}

// RegisterHandler gère la requête d'inscription d'un utilisateur
// @Summary Inscription d'un utilisateur
// @Description Inscription d'un nouvel utilisateur
// @Tags Auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "Informations de l'utilisateur"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/register [post]
func (ctrl *AuthController) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if req.Role == "" {
		req.Role = "Player"
	}

	userInfo := models.Users{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Location:     req.Location,

		ConfirmationToken: ulid.MustNew(ulid.Timestamp(time.Now()), ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)).String(),
	}

	user, err := ctrl.AuthService.RegisterUser(userInfo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// Envoyer un email de confirmation
	err = ctrl.AuthService.EmailService.SendConfirmationEmail(user.Email, user.ConfirmationToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send confirmation email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Inscription réussie, veuillez vérifier votre email pour confirmer votre compte"})
}

// RegisterRequest représente le corps de la requête d'inscription
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
	Location string `json:"location"`
	Role     string `json:"role"`
}

// UserResponse représente la réponse de l'inscription d'un utilisateur
type UserResponse struct {
	Username          string `json:"username"`
	Email             string `json:"email"`
	Location          string `json:"location"`
	Role              string `json:"role"`
	ConfirmationToken string `json:"confirmationToken"`
}

// UserHandler gère la requête pour récupérer les informations de l'utilisateur connecté
// @Summary Récupérer les informations de l'utilisateur connecté
// @Description Récupérer les informations de l'utilisateur connecté
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} models.Users
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/userInfo [get]
func (ctrl *AuthController) UserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	user, err := ctrl.AuthService.GetUserByID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UserUpdate gère la requête pour mettre à jour les informations de l'utilisateur connecté
// @Summary Mettre à jour les informations de l'utilisateur connecté
// @Description Mettre à jour les informations de l'utilisateur connecté
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param username body string false "Nom d'utilisateur"
// @Param email body string false "Adresse email"
// @Param password body string false "Mot de passe"
// @Param profilePhoto body string false "Photo de profil"
// @Param favoriteSport body string false "Sport favori"
// @Param bio body string false "Biographie"
// @Param location body string false "Localisation"
// @Param birthDate body string false "Date de naissance (YYYY-MM-DD)"
// @Param role body string false "Rôle"
// @Param skillLevel body string false "Niveau de compétence"
// @Param pac body int false "Pac"
// @Param sho body int false "Sho"
// @Param pas body int false "Pas"
// @Param dri body int false "Dri"
// @Param def body int false "Def"
// @Success 200 {object} models.Users
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/userUpdate [put]
func (ctrl *AuthController) UserUpdate(c *fiber.Ctx) error {
	userID := c.Locals("user_id")
	if userID == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userIDStr, ok := userID.(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Invalid user ID"})
	}

	var req struct {
		Username      string `json:"username"`
		Email         string `json:"email"`
		Password      string `json:"password"`
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

	var birthDate *time.Time
	if req.BirthDate != "" {
		parsedDate, err := time.Parse("2006-01-02", req.BirthDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid birth date format. Use YYYY-MM-DD"})
		}
		birthDate = &parsedDate
	}

	role := models.Role(req.Role)

	user, err := ctrl.AuthService.UpdateUser(
		userIDStr, req.Username, req.Email, req.Password, req.ProfilePhoto, req.FavoriteSport,
		req.Location, req.Bio, birthDate, role, req.SkillLevel, req.Pac, req.Sho, req.Pas,
		req.Dri, req.Def, req.Phy, req.MatchesPlayed, req.MatchesWon, req.GoalsScored, req.BehaviorScore,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// GoogleLogin redirige l'utilisateur vers la page de connexion Google
// @Summary Rediriger l'utilisateur vers la page de connexion Google
// @Description Rediriger l'utilisateur vers la page de connexion Google
// @Tags Auth
// @Produce json
// @Success 302 {string} string
// @Router /api/auth/google [get]
// GoogleLogin redirige l'utilisateur vers la page de connexion Google
// GoogleLogin redirige l'utilisateur vers la page de connexion Google
func (ctrl *AuthController) GoogleLogin(c *fiber.Ctx) error {
	url := ctrl.AuthService.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

// GoogleCallback gère le callback de Google après l'authentification
func (ctrl *AuthController) GoogleCallback(c *fiber.Ctx) error {
	var req struct {
		IDToken string `json:"idToken"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	token, err := ctrl.AuthService.GoogleOauthConfig.TokenSource(context.Background(), &oauth2.Token{
		AccessToken: req.IDToken,
	}).Token()
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	client := ctrl.AuthService.GoogleOauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get user info"})
	}
	defer resp.Body.Close()

	var userInfo struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode user info"})
	}

	// Rechercher l'utilisateur dans la base de données
	user, err := ctrl.AuthService.GetUserByEmail(userInfo.Email)
	if err != nil {
		// Si l'utilisateur n'existe pas, créez un nouvel utilisateur
		user, err = ctrl.AuthService.RegisterUser(models.Users{
			Email: userInfo.Email,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to register user"})
		}
	}

	// Générer un token JWT pour l'utilisateur
	userID, err := ulid.Parse(user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to parse user ID"})
	}
	accessToken, err := middlewares.GenerateToken(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate access token"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": accessToken})
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginHandler gère la requête de connexion d'un utilisateur
// @Summary Connexion d'un utilisateur
// @Description Connexion d'un utilisateur
// @Tags Auth
// @Accept json
// @Produce json
// @Param login body LoginRequest true "Email et mot de passe"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/login [post]
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
// @Summary Rafraîchir le token d'un utilisateur
// @Description Rafraîchir le token d'un utilisateur
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshToken body string true "Refresh Token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/refresh [post]
func (ctrl *AuthController) RefreshHandler(c *fiber.Ctx) error {
	var req struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newAccessToken, newRefreshToken, err := ctrl.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

func (ctrl *AuthController) ConfirmEmailHandler(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Token manquant"})
	}

	var user models.Users
	if err := ctrl.AuthService.DB.Where("confirmation_token = ?", token).First(&user).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Utilisateur non trouvé"})
	}

	user.IsConfirmed = true
	user.ConfirmationToken = ""

	if err := ctrl.AuthService.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Erreur lors de la confirmation de l'email"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Email confirmé avec succès"})
}

// DeleteUserHandler gère la demande de suppression d'un utilisateur

// @Summary Supprimer un utilisateur
// @Description Supprimer un utilisateur
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/deleteMyAccount [delete]
func (ctrl *AuthController) DeleteUserHandler(c *fiber.Ctx) error {
	userID, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user ID"})
	}
	matchId, ok := c.Locals("match_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid match ID"})
	}

	err := ctrl.AuthService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = ctrl.MatchService.DeleteMatch(userID, matchId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Utilisateur supprimé avec succès"})
}

// GetUsersHandler gère la demande de récupération de tous les utilisateurs
// @Summary Récupérer tous les utilisateurs
// @Description Récupérer tous les utilisateurs
// @Tags Auth
// @Produce json
// @Success 200 {object} []models.Users
// @Failure 500 {object} map[string]interface{}
// @Router /api/users [get]
func (ctrl *AuthController) GetUsersHandler(c *fiber.Ctx) error {
	users, err := ctrl.AuthService.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(users)
}

//public profile by id

// @Summary GetPublicUserInfoHandler
// @Description Get public user info by id
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} models.Users
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{id}/public [get]
func (ctrl *AuthController) GetPublicUserInfoHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	publicInfo, err := ctrl.AuthService.GetPublicUserInfoByID(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(publicInfo)
}

// @Summary Attribuer le rôle d'arbitre à un joueur
// @Description Attribuer le rôle d'arbitre à un joueur
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Param organizerID path string true "Organizer ID"
// @Param playerID path string true "Player ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
func (ctrl *AuthController) AssignRefereeRole(c *fiber.Ctx) error {
	organizerID := c.Params("organizerID")
	playerID := c.Params("playerID")

	err := ctrl.AuthService.AssignRefereeRole(organizerID, playerID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Role assigned successfully"})
}

// @Summary Mettre à jour les statistiques d'un joueur
// @Description Mettre à jour les statistiques d'un joueur
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
func (ctrl *AuthController) UpdateUserStatistics(c *fiber.Ctx) error {
	userID := c.Params("id")
	var stats struct {
		MatchesPlayed int `json:"matchesPlayed"`
		MatchesWon    int `json:"matchesWon"`
		GoalsScored   int `json:"goalsScored"`
		BehaviorScore int `json:"behaviorScore"`
	}

	if err := c.BodyParser(&stats); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	err := ctrl.AuthService.UpdateUserStatistics(userID, stats.MatchesPlayed, stats.MatchesWon, stats.GoalsScored, stats.BehaviorScore)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User statistics updated successfully"})
}
