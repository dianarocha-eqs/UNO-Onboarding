package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/util"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Interface for user's services
type UserService interface {
	// Creates a new user and returns the user's UUID
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	// Updates an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Retrieves a user by their ID
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
}

// Handles user's logic and interaction with the repository
type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

// Checks the required fields of the user
func validateRquiredFields(user *domain.User) error {
	if user.Email == "" || user.Name == "" || user.Phone == "" {
		return errors.New("name, email, and phone are required fields")
	}
	return nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := validateRquiredFields(user); err != nil {
		return "", err
	}
	// Generate UUID
	user.ID = uuid.New().String() // Generate UUID as string

	// Generate random password and hash
	plainPassword, hashedPassword, err := util.GenerateRandomPasswordHash()
	if err != nil {
		return "", err
	}
	// Assign the hashed password to the user
	user.Password = hashedPassword

	err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}

	// Send the plain password to the user's email
	emailSubject := "Welcome to UNO Service"
	emailBody := fmt.Sprintf("Hello %s,\n\nYour account has been created. Your temporary password is: %s\n\nPlease change it after logging in.", user.Name, plainPassword)

	err = util.SendEmail(user.Email, emailSubject, emailBody)
	if err != nil {
		return "", errors.New("user created but failed to send email")
	}
	return user.ID, nil
}

// Updates an existing user
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	// Validate required fields
	if err := validateRquiredFields(user); err != nil {
		return err
	}

	// Attempt to update user in the database
	if err := s.Repo.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", user.ID, err)
	}

	return nil
}

// Retrieves a user by their ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.Repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with id %s not found", id)
		}
		return nil, fmt.Errorf("error retrieving user with id %s: %v", id, err)
	}
	return user, nil
}
