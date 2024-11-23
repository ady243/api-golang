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


// GenerateToken génère un nouveau token JWT pour un utilisateur donné
//
// Le token est valable pour 24h.
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


// ParseToken décode un token JWT et retourne les claims associés.
//
// Cette fonction prend un token JWT (tokenString) en entrée et utilise la clé secrète pour le décoder.
// Si le token est valide, les claims (informations contenues dans le token) sont retournés.
// En cas d'erreur, elle retourne une erreur appropriée, comme une clé secrète manquante,
// un token invalide ou un problème lors du parsing.
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


// GenerateRefreshToken generates a new refresh token for a given user.
// The token is valid for 72 hours.
// It requires a valid SECRET_KEY environment variable.
// Returns the signed token as a string or an error if the secret key is missing.
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



// JWTMiddleware is a middleware that checks for a valid JWT token in the Authorization header of the request.
// If the token is valid, it extracts the user ID and role from the token and stores them in the Locals of the request.
// It also extracts the permissions for the given role and stores them in the Locals.
// If the token is invalid or missing, it returns a 401 status code with an appropriate error message.
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