package repository

import (
	config "api/configs"
	auth_domain "api/internal/auth/domain"
	"context"
	"fmt"

	uuid "github.com/tentone/mssql-uuid"
	"gorm.io/gorm"
)

// Interface for token's data operations
type AuthRepository interface {
	// Stores a new token in the database
	StoreToken(ctx context.Context, auth *auth_domain.AuthToken) error
	// Gets a specific token and returns it
	GetToken(ctx context.Context, tokenStr string) (*auth_domain.AuthToken, error)
	// Sets token to false (invalid)
	InvalidateToken(ctx context.Context, tokenStr string) error
	// Gets last valid token from user and returns it
	GetTokenByUserID(ctx context.Context, userID uuid.UUID) (*auth_domain.AuthToken, error)
}

// GORM to interact with the token's database
type AuthRepositoryImpl struct {
	DB *gorm.DB
}

// Connects with the database
func NewAuthRepository() (AuthRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &AuthRepositoryImpl{DB: db}, nil
}

func (r *AuthRepositoryImpl) StoreToken(ctx context.Context, auth *auth_domain.AuthToken) error {
	return r.DB.WithContext(ctx).Create(auth).Error
}

func (r *AuthRepositoryImpl) GetToken(ctx context.Context, tokenStr string) (*auth_domain.AuthToken, error) {
	var auth auth_domain.AuthToken
	err := r.DB.WithContext(ctx).Where("token = ?", tokenStr).First(&auth).Error
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

func (r *AuthRepositoryImpl) InvalidateToken(ctx context.Context, tokenStr string) error {
	return r.DB.WithContext(ctx).Model(&auth_domain.AuthToken{}).
		Where("token = ?", tokenStr).
		Update("is_valid", false).Error
}

func (r *AuthRepositoryImpl) GetTokenByUserID(ctx context.Context, userID uuid.UUID) (*auth_domain.AuthToken, error) {
	var authToken auth_domain.AuthToken

	err := r.DB.WithContext(ctx).
		Where("user_id = ? AND is_valid = ?", userID, true).
		First(&authToken).Error

	if err != nil {
		return nil, err
	}

	return &authToken, nil
}
