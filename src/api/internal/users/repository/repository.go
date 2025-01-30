package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"fmt"

	uuid "github.com/tentone/mssql-uuid"

	"gorm.io/gorm"
)

// Interface for user's data operations
type UserRepository interface {
	// Stores a new user in the database
	CreateUser(ctx context.Context, user *domain.User) error
	// Updates the details of an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get the users info
	GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error)
	// Get any users info
	GetUsers(ctx context.Context) ([]domain.User, error)
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
	return r.DB.WithContext(ctx).Model(&domain.User{}).Where("id = ?", user.ID).Updates(user).Error
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.DB.WithContext(ctx).Where("id = ?", userID).First(&user).Error
	fmt.Printf("%v", user.Role)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) GetUsers(ctx context.Context) ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Find(&users).Error
	return users, err
}
