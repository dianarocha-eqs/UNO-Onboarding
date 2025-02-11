package util

import (
	"api/configs"
	"api/internal/users/domain"
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/tentone/mssql-uuid"
)

var jwtSecret []byte

// init function to ensure JWT secret is set from config.json
func init() {
	// Load configuration from configs/config.json
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Ensure JWT secret is set
	if config.JWTSecret == "" {
		log.Fatal("JWT secret is missing in config.json")
	}

	// Store JWT secret
	jwtSecret = []byte(config.JWTSecret)
}

type JWT interface {
	GenerateJWT(user *domain.User) (string, error)
	ValidateJWT(tokenString string) (*Claims, error)
	GenerateSecureToken() (string, error)
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate a new JWT token
func GenerateJWT(user *domain.User) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validate JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
