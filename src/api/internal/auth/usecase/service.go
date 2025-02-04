package usecase

import (
	auth_domain "api/internal/auth/domain"
	auth_repos "api/internal/auth/repository"
	jwt "api/internal/auth/util"
	user_domain "api/internal/users/domain"
	user_repos "api/internal/users/repository"
	"context"
	"errors"
	"fmt"
)

// Interface for authentication services
type AuthService interface {
	// Generates the token and stores it
	AddToken(ctx context.Context, user *user_domain.User) (string, error)
	// Validate the token and removes it
	RemoveToken(ctx context.Context, tokenStr string) error
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

	// Check if the user already has a valid token (which means is already logged in)
	existingToken, err := s.AuthRepo.GetTokenByUserID(ctx, user.ID)
	if err == nil && existingToken.IsValid {
		return "", errors.New("user already logged in with an active session")
	}

	// Generate JWT token
	tokenStr, err := jwt.GenerateJWT(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	// Store token data
	authToken := &auth_domain.AuthToken{
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

func (s *AuthServiceImpl) RemoveToken(ctx context.Context, tokenStr string) error {
	// Validate the token first
	_, err := jwt.ValidateJWT(tokenStr)
	if err != nil {
		return errors.New("invalid token, cannot logout")
	}

	// Retrieve the token from the database
	authToken, err := s.AuthRepo.GetToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("failed to get token: %v", err)
	}

	// If the token is already invalidated, it means it's already disconnected
	if !authToken.IsValid {
		return errors.New("token is already invalidated")
	}

	// sets token validation to false
	err = s.AuthRepo.InvalidateToken(ctx, tokenStr)
	if err != nil {
		return fmt.Errorf("failed to invalidate token: %v", err)
	}

	return nil
}

func (s *AuthServiceImpl) IsTokenValid(ctx context.Context, tokenStr string) (bool, error) {
	// Retrieve the token from the database
	authToken, err := s.AuthRepo.GetToken(ctx, tokenStr)
	if err != nil {
		return false, fmt.Errorf("failed to get token: %v", err)
	}

	// Check if the token is valid
	return authToken.IsValid, nil
}
