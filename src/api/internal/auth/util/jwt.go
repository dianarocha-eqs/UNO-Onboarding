package util

import (
	"api/internal/users/domain"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	uuid "github.com/tentone/mssql-uuid"
)

var jwtSecret []byte

func init() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Retrieve the JWT secret key
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	jwtSecret = []byte(secret)
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
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 72 hours

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
