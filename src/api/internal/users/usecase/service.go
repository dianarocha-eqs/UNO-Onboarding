package usecase

import (
	aux "api/auxiliary"
	auth_repository "api/internal/auth/repository"
	"api/internal/users/domain"
	users_repository "api/internal/users/repository"
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
	// Reset previous password of user
	ResetPassword(ctx context.Context, token string, newPassword string) error
}

// Handles user's logic and interaction with the repository
type UserServiceImpl struct {
	UserRepository users_repository.UserRepository
	AuthRepository auth_repository.AuthRepository
}

func NewUserService(userRepo users_repository.UserRepository, authRepo auth_repository.AuthRepository) UserService {
	return &UserServiceImpl{
		UserRepository: userRepo,
		AuthRepository: authRepo,
	}
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
	err = s.UserRepository.CreateUser(ctx, user)
	if err != nil {
		return uuid.NilUUID, errors.New("failed to create user")
	}

	// Send the plain password to the user's email
	var emailSubject string
	var emailBody string
	emailSubject = "Welcome to UNO Service"
	emailBody = fmt.Sprintf("Hello %s,\n\nYour account has been created. Your temporary password is: %s\n\nPlease change it after logging in.", user.Name, plainPasswordForEmail)

	err = utils.CreateEmail(user.Email, emailSubject, emailBody)
	if err != nil {
		return user.ID, errors.New("user created but failed to send email")
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

	err = s.UserRepository.UpdateUser(ctx, user)
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
	users, err = s.UserRepository.ListUsers(ctx, search, sortDirection)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve users: %w", err)
	}

	if search != "" && len(users) == 0 {
		return nil, errors.New("no result was found")
	}

	return users, nil
}

func (s *UserServiceImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error) {
	user, err := s.UserRepository.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}

func (s *UserServiceImpl) GetRoutesAuthorization(ctx context.Context, tokenStr string, getRole *bool, getUserID *uuid.UUID) error {
	err := s.UserRepository.GetRoutesAuthorization(ctx, tokenStr, getRole, getUserID)
	if err != nil {
		return fmt.Errorf("failed to retrieve user id: %v", err)
	}
	return err
}

func (s *UserServiceImpl) ResetPassword(ctx context.Context, token string, newPassword string) error {

	// hash the password received to store in database (never plain password)
	_, hashedPassword, err := utils.GeneratePasswordHash(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = s.UserRepository.ResetPassword(ctx, token, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update user's password: %w", err)
	}

	return nil
}
