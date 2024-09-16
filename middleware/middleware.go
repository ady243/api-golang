package middlewares

import (
    "errors"
    "github.com/dgrijalva/jwt-go"
    "net/http"
    "os"
    "context"
	"time"
)


type Claims struct {
    UserID uint `json:"user_id"`
    jwt.StandardClaims
}

// GenerateToken génère un nouveau JWT pour un utilisateur donné
func GenerateToken(userID uint) (string, error) {
    secretKey := os.Getenv("SECRET_KEY")
    if secretKey == "" {
        return "", errors.New("SECRET_KEY not found")
    }
    claims := Claims{
        UserID: userID,
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

// JWTMiddleware vérifie le token JWT dans les requêtes HTTP
func JWTMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tokenString := r.Header.Get("Authorization")
        if tokenString == "" {
            http.Error(w, "Missing token", http.StatusUnauthorized)
            return
        }

        claims, err := ParseToken(tokenString)
        if err != nil {
            http.Error(w, "Invalid token", http.StatusUnauthorized)
            return
        }

        // Ajoute les claims au contexte de la requête
        ctx := context.WithValue(r.Context(), "claims", claims)
        r = r.WithContext(ctx)

        next.ServeHTTP(w, r)
    })
}
