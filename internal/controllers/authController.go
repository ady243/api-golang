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

// AuthController fournit les gestionnaires pour les opérations d'authentification
type AuthController struct {
    AuthService *services.AuthService
    ImageService *services.ImageService
}

func NewAuthController(authService *services.AuthService, imageService *services.ImageService) *AuthController {
    return &AuthController{
        AuthService:  authService,
        ImageService: imageService,
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


// UserHandler gère la requête pour récupérer les informations de l'utilisateur connecté
func (ctrl *AuthController) UserHandler(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string) 

    user, err := ctrl.AuthService.GetUserByID(userID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(user)
}
func (ctrl *AuthController) UserUpdate(c *fiber.Ctx) error {
    userID := c.Locals("userID").(string)
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
        userID, req.Username, req.Email, req.Password, req.ProfilePhoto, req.FavoriteSport,
        req.Location, req.Bio, birthDate, role, req.SkillLevel, req.Pac, req.Sho, req.Pas,
        req.Dri, req.Def, req.Phy, req.MatchesPlayed, req.MatchesWon, req.GoalsScored, req.BehaviorScore,
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(user)
}
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
        user, err = ctrl.AuthService.Register("", userInfo.Email, "", "", "", "", "", nil, models.Role("Player"), "", 0, 0, 0, 0, 0, 0, 0, 0, 0, 0)
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