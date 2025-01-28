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
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
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
	return r.DB.WithContext(ctx).Save(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, err
	}
	return &user, nil
}
