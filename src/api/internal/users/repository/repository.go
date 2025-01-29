package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

// Interface for user's data operations
type UserRepository interface {
	// Stores a new user in the database
	CreateUser(ctx context.Context, user *domain.User) error
	// Updates the details of an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get the user's info
	GetUserByID(ctx context.Context, userID string) (*domain.User, error)
}

// Performs user's data operations using GORM to interact with the database
type UserRepositoryImpl struct {
	DB *gorm.DB
}

// Connects with the database
func NewUserRepository() (UserRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &UserRepositoryImpl{DB: db}, nil
}

func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

func (r *UserRepositoryImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Select("name", "email", "phone", "password", "picture").Updates(user).Error
}

// This was already created on this branch mainly for password change
func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
