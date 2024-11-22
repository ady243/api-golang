package controllers

import (
	"context"
	"encoding/json"
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
}

// NewAuthController creates a new instance of AuthController.
// It requires an AuthService and an ImageService to handle authentication
// and image-related operations, respectively.
func NewAuthController(authService *services.AuthService, imageService *services.ImageService) *AuthController {
	return &AuthController{
		AuthService:  authService,
		ImageService: imageService,
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
		return c.Status(fiber.StatusBadRequest).JSON(map[string]interface{}{"error": err.Error()})
	}

	if req.Role == "" {
		req.Role = "Player"
	}
	role := models.Role(req.Role)

	userInfo := models.Users{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: req.Password,
		Location:     req.Location,
		Role:         role,
	}

	user, err := ctrl.AuthService.RegisterUser(userInfo)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]interface{}{"error": err.Error()})
	}

	// Créer une réponse personnalisée
	userResponse := UserResponse{
		Username: user.Username,
		Email:    user.Email,
		Location: user.Location,
		Role:     user.Role,
	}

	return c.Status(fiber.StatusOK).JSON(userResponse)
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
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Location string      `json:"location"`
	Role     models.Role `json:"role"`
}

// UserHandler gère la requête pour récupérer les informations de l'utilisateur connecté
// @Summary Récupérer les informations de l'utilisateur connecté
// @Description Récupérer les informations de l'utilisateur connecté
// @Tags Auth
// @Security Bearer
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
// @Security Bearer
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
	userID := c.Locals("userID")
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

// GetPublicUserInfoHandler gère la requête pour récupérer les informations de l'utilisateur public
// @Summary Récupérer les informations de l'utilisateur public
// @Description Récupérer les informations de l'utilisateur public
// @Tags Auth
// @Security Bearer
// @Produce json
// @Param id path string true "ID de l'utilisateur"
// @Success 200 {object} models.Users
// @Failure 404 {object} map[string]interface{}
// @Router /api/users/{id}/public [get]
func (ctrl *AuthController) GetPublicUserInfoHandler(c *fiber.Ctx) error {
	userID := c.Params("id")
	publicInfo, err := ctrl.AuthService.GetPublicUserInfo(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}

	return c.JSON(publicInfo)
}

// GoogleLogin redirige l'utilisateur vers la page de connexion Google
// @Summary Rediriger l'utilisateur vers la page de connexion Google
// @Description Rediriger l'utilisateur vers la page de connexion Google
// @Tags Auth
// @Produce json
// @Success 302 {string} string
// @Router /api/auth/google [get]
func (ctrl *AuthController) GoogleLogin(c *fiber.Ctx) error {
	url := ctrl.AuthService.GoogleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

// GoogleCallback gère le callback de Google après l'authentification
// @Summary Gérer le callback de Google après l'authentification
// @Description Gérer le callback de Google après l'authentification
// @Tags Auth
// @Accept json
// @Produce json
// @Param idToken body string true "ID Token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/auth/google/callback [get]
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
			Role:  models.Role("Player"),
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
	accessToken, err := middlewares.GenerateToken(userID, user.Role)
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

	newAccessToken, err := ctrl.AuthService.Refresh(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"accessToken": newAccessToken})
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
// @Security Bearer
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

	err := ctrl.AuthService.DeleteUser(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Utilisateur supprimé avec succès"})
}
