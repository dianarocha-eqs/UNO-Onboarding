package usecase

import (
	aux "api/auxiliary"
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/utils"

	"context"
	"errors"
	"fmt"
	"strings"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for user's services
type UserService interface {
	// Creates a new user and returns the user's UUID
	CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error)
	// Updates an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get users
	ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error)
	// Get user by email and password
	GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error)
	// Checks user's role and uuid from token
	GetRoutesAuthorization(ctx context.Context, tokenStr string, getRole *bool, getUserID *uuid.UUID) error
	// Only update password
	UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error
}

// Handles user's logic and interaction with the repository
type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

// Checks the required fields and their format
func validateRequiredFields(user *domain.User) error {
	// trim spaces for required fields only
	user.Name = strings.Join(strings.Fields(user.Name), " ")
	user.Email = strings.Join(strings.Fields(user.Email), " ")
	user.Phone = strings.Join(strings.Fields(user.Phone), " ")

	if user.Name == "" || user.Email == "" || user.Phone == "" {
		return errors.New("name, email, and phone are required fields")
	}

	var err error
	// Validate email
	if err = aux.ValidateEmail(user.Email); err != nil {
		return err
	}
	// Validate phone number
	if err = aux.ValidatePhone(user.Phone); err != nil {
		return err
	}

	return nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	var err error
	if err = validateRequiredFields(user); err != nil {
		return uuid.NilUUID, err
	}

	// Create the hashed password
	var plainPasswordForEmail string
	plainPasswordForEmail, user.Password, err = utils.GeneratePasswordHash("")
	if err != nil {
		return uuid.NilUUID, err
	}

	user.ID = uuid.NewV4()
	err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.NilUUID, errors.New("failed to create user")
	}

	// Send the email with the plain password (only after user is created)
	if err = utils.SendEmail(user, plainPasswordForEmail); err != nil {
		return uuid.UUID{}, err
	}

	return user.ID, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID == uuid.NilUUID {
		return errors.New("user ID is required")
	}

	var err error
	if err = validateRequiredFields(user); err != nil {
		return err
	}

	if user.Password != "" {
		var err error
		_, user.Password, err = utils.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}
	}

	err = s.Repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", user.ID.String(), err)
	}

	return nil
}

func (s *UserServiceImpl) ListUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error) {
	if sortDirection != 1 && sortDirection != -1 && sortDirection != 0 {
		return nil, errors.New("invalid sort direction: must be 1 (ASC) or -1 (DESC) or 0 (NO ORDER)")
	}

	var users []domain.User
	var err error
	// Call the repository (which executes the stored procedure to handle searching and sorting)
	users, err = s.Repo.ListUsers(ctx, search, sortDirection)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	if search != "" && len(users) == 0 {
		return nil, errors.New("no result was found")
	}

	return users, nil
}

func (s *UserServiceImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.Repo.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}

func (s *UserServiceImpl) GetRoutesAuthorization(ctx context.Context, tokenStr string, getRole *bool, getUserID *uuid.UUID) error {
	return s.Repo.GetRoutesAuthorization(ctx, tokenStr, getRole, getUserID)
}

func (s *UserServiceImpl) UpdatePassword(ctx context.Context, userID uuid.UUID, password string) error {

	var err error
	err = s.Repo.UpdatePassword(ctx, userID, password)
	if err != nil {
		return fmt.Errorf("failed to update password for user with id %s: %v", userID, err)
	}

	return nil
}
