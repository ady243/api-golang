package middlewares

import (
    "errors"
    "os"
    "strings"
    "time"

    "github.com/ady243/teamup/helpers"
    "github.com/ady243/teamup/internal/models"
    "github.com/dgrijalva/jwt-go"
    "github.com/gofiber/fiber/v2"
    "github.com/oklog/ulid/v2"
)

type Claims struct {
    UserID ulid.ULID `json:"user_id"`
    Role   models.Role `json:"role"`
    jwt.StandardClaims
}

// GenerateToken génère un nouveau JWT pour un utilisateur donné
func GenerateToken(userID ulid.ULID, role models.Role) (string, error) {
    secretKey := os.Getenv("SECRET_KEY")
    if secretKey == "" {
        return "", errors.New("SECRET_KEY not found")
    }
    claims := Claims{
        UserID: userID,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// ParseToken vérifie le JWT et retourne les claims
func ParseToken(tokenString string) (*Claims, error) {
    secretKey := os.Getenv("SECRET_KEY")
    if secretKey == "" {
        return nil, errors.New("SECRET_KEY not found")
    }
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(secretKey), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// GenerateRefreshToken génère un nouveau refreshToken pour un utilisateur donné
func GenerateRefreshToken(userID ulid.ULID, role models.Role) (string, error) {
    secretKey := os.Getenv("SECRET_KEY")
    if secretKey == "" {
        return "", errors.New("SECRET_KEY not found")
    }
    claims := Claims{
		UserID: userID,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
            IssuedAt:  time.Now().Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secretKey))
}

// JWTMiddleware vérifie le token JWT dans les requêtes HTTP et ajoute les permissions de l'utilisateur
func JWTMiddleware(c *fiber.Ctx) error {
    authHeader := c.Get("Authorization")
    if authHeader == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing token"})
    }

    authParts := strings.Split(authHeader, " ")
    if len(authParts) != 2 {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
    }
    tokenString := authParts[1]

    claims, err := ParseToken(tokenString)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
    }

    permissions := helpers.GetPermissions(claims.Role)

    c.Locals("user_id", claims.UserID.String())
    c.Locals("user_role", claims.Role)
    c.Locals("permissions", permissions)
    return c.Next()
}