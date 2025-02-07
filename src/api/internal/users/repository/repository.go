package repository

import (
	config "api/configs"
	"api/internal/users/domain"
	"context"
	"database/sql"
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
	ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error)
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
	// any of the fields will sent null if empty, except picture
	// procedure receives the null and doesnt change those fields
	return r.DB.Exec(
		"EXEC UpdateUser ?, ?, ?, ?, ?, ?",
		user.ID,
		user.Name,
		user.Email,
		user.Phone,
		user.Picture,
		sql.NullString{String: user.Password, Valid: user.Password != ""},
	).Error
}

func (r *UserRepositoryImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.DB.Raw("EXEC GetUserByID ?", userID).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryImpl) ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error) {
	var users []domain.User
	err := r.DB.Raw("EXEC ListUsers ?, ?", search, sortDirection).Scan(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}
