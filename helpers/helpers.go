package helpers

import (
    "golang.org/x/crypto/bcrypt"
    "log"
)

// HashPassword prend un mot de passe en clair et le hache en utilisant bcrypt
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Printf("Error hashing password: %v", err)
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPasswordHash vérifie que le mot de passe fourni correspond au hash
func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        log.Printf("Password mismatch or hash error: %v", err)
        return false
    }
    return true
}

