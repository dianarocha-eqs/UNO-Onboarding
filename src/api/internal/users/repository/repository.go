package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindUserByRole(ctx context.Context, role bool) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	// GetAllUsers(ctx context.Context) ([]domain.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository() (UserRepository, error) {
	db, err := config.ConnectDB()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	return &UserRepositoryImpl{DB: db}, nil
}

func (r *UserRepositoryImpl) FindUserByRole(ctx context.Context, role bool) (*domain.User, error) {
	var user domain.User
	result := r.DB.WithContext(ctx).Where("role = ?", role).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}
func (r *UserRepositoryImpl) CreateUser(ctx context.Context, user *domain.User) error {
	return r.DB.WithContext(ctx).Create(user).Error
}

// func (r *UserRepositoryImpl) GetAllUsers() ([]domain.User, error) {
// 	var users []domain.User
// 	err := r.DB.Find(&users).Error
// 	return users, err
// }
