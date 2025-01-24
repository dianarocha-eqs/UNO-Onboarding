package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"fmt"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
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
func (r *UserRepositoryImpl) CreateUser(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepositoryImpl) GetAllUsers() ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}
