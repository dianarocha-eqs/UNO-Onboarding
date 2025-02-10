package usecase

import (
	auth_domain "api/internal/auth/domain"
	auth_repos "api/internal/auth/repository"
	jwt "api/internal/auth/util"
	user_domain "api/internal/users/domain"
	user_repos "api/internal/users/repository"
	"context"
	"fmt"
)

// Interface for authentication services
type AuthService interface {
	// Generates the token and stores it
	AddToken(ctx context.Context, user *user_domain.User) (string, error)
	// Sets token to invalid
	InvalidateToken(ctx context.Context, tokenStr string) error
	// Checks the state of token
	IsTokenValid(ctx context.Context, tokenStr string) (bool, error)
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

	var tokenStr string
	var err error
	// Generate JWT token
	tokenStr, err = jwt.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	var authToken *auth_domain.AuthToken
	// Store token data
	authToken = &auth_domain.AuthToken{
		UserID:  user.ID,
		Token:   tokenStr,
		IsValid: true,
	}

	// in database
	err = s.AuthRepo.StoreToken(ctx, authToken)
	if err != nil {
		return "", fmt.Errorf("failed to store token: %v", err)
	}

	return tokenStr, nil
}

func (s *AuthServiceImpl) InvalidateToken(ctx context.Context, tokenStr string) error {
	var err error
	// sets token validation to false
	err = s.AuthRepo.InvalidateToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	return nil
}

func (s *AuthServiceImpl) IsTokenValid(ctx context.Context, tokenStr string) (bool, error) {
	// Validate the token first
	_, err := jwt.ValidateJWT(tokenStr)
	if err != nil {
		return false, fmt.Errorf("invalid JWT: %v", err)
	}

	// Retrieve the token from the database
	authToken, err := s.AuthRepo.GetToken(ctx, tokenStr)
	if err != nil {
		return false, fmt.Errorf("failed to get token: %v", err)
	}

	return authToken.IsValid, nil
}
