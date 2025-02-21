package usecase

import (
	auth_domain "api/internal/auth/domain"
	auth_repos "api/internal/auth/repository"
	jwt "api/internal/auth/util"
	user_domain "api/internal/users/domain"
	user_repos "api/internal/users/repository"
	"context"
	"fmt"
	"time"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for authentication services
type AuthService interface {
	// Generates the token and stores it
	AddToken(ctx context.Context, user *user_domain.User) (string, error)
	// Sets token to invalid
	InvalidateToken(ctx context.Context, tokenStr string) error
	// Checks the state of token
	IsTokenValid(ctx context.Context, tokenStr string) (bool, error)
	// Generates the token for password recovery
	AddTokenForPasswordRecovery(ctx context.Context, user *user_domain.User) (string, error)
}

type AuthServiceImpl struct {
	AuthRepo auth_repos.AuthRepository
	UserRepo user_repos.UserRepository
}

func NewAuthService(authRepo auth_repos.AuthRepository, userRepo user_repos.UserRepository) AuthService {
	return &AuthServiceImpl{
		AuthRepo: authRepo,
		UserRepo: userRepo,
	}
}

func (s *AuthServiceImpl) AddToken(ctx context.Context, user *user_domain.User) (string, error) {
	// Generate JWT token
	var tokenStr, err = jwt.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Parse the token and retrieve the expiration time from the JWT claims
	claims, err := jwt.ValidateJWT(tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %v", err)
	}

	// Store token data
	var authToken = &auth_domain.AuthToken{
		ID:      uuid.NewV4(),
		UserID:  user.ID,
		Token:   tokenStr,
		IsValid: true,
		// Set the created time
		CreatedAt: time.Now().UTC(),
		// Set the expiration time (same as jwt token)
		ExpiredAt: claims.ExpiresAt.Time,
	}

	// in database
	err = s.AuthRepo.StoreToken(ctx, authToken)
	if err != nil {
		return "", fmt.Errorf("failed to store token: %v", err)
	}

	return tokenStr, nil
}

func (s *AuthServiceImpl) InvalidateToken(ctx context.Context, tokenStr string) error {
	// Sets token validation to false
	var err = s.AuthRepo.InvalidateToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	return nil
}

func (s *AuthServiceImpl) IsTokenValid(ctx context.Context, tokenStr string) (bool, error) {
	// Validate the token first (structure and expiration)
	var _, err = jwt.ValidateJWT(tokenStr)
	if err != nil {
		fmt.Print(err)
		return false, fmt.Errorf("invalid or expired JWT: %v", err)
	}

	var authToken *auth_domain.AuthToken
	// Retrieve the token from the database
	authToken, err = s.AuthRepo.GetToken(ctx, tokenStr)
	if err != nil {
		return false, fmt.Errorf("failed to get token: %v", err)
	}

	// Returns the token state (valid or not)
	return authToken.IsValid, nil
}

func (s *AuthServiceImpl) AddTokenForPasswordRecovery(ctx context.Context, user *user_domain.User) (string, error) {
	// Generate a JWT token for password recovery
	tokenStr, err := jwt.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Parse the token and retrieve the expiration time from the JWT claims
	claims, err := jwt.ValidateJWT(tokenStr)
	if err != nil {
		return "", fmt.Errorf("failed to validate token: %v", err)
	}

	var expirationTime = time.Now().UTC().Add(10 * time.Minute)
	claims.ExpiresAt.Time = expirationTime

	// Insert the token data into the Users_Tokens table for password recovery
	err = s.AuthRepo.StoreTokenToPasswordRecovery(ctx, user.ID, tokenStr, expirationTime)
	if err != nil {
		return "", fmt.Errorf("failed to store token for password recovery: %v", err)
	}

	return tokenStr, nil
}
