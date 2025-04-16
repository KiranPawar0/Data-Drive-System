package jwthelper

import (
	"fmt"
	"time"

	env "github.com/KiranPawar0/Data-Drive-System/pkg/config"
	"github.com/golang-jwt/jwt"
)

// GenerateJWTToken generates a JWT token with the provided user ID
func GenerateJWTToken(email string) (string, error) {
	// Load environment config
	cfg, err := env.Env()
	if err != nil {
		return "", fmt.Errorf("failed to load config: %v", err)
	}

	var jwtSecret = []byte(cfg.JWT.Secret)

	// Create a new token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set token claims
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = time.Now().Add(360 * 24 * time.Hour).Unix() // Token expiration time

	// Sign the token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
