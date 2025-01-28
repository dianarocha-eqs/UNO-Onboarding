package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/utils"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

// Interface for user's services
type UserService interface {
	// Creates a new user and returns the user's UUID
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	// Updates an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
}

// Handles user's logic and interaction with the repository
type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

// Checks the required fields of the user
func validateRequiredFields(user *domain.User) error {
	if user.Email == "" || user.Name == "" || user.Phone == "" {
		return errors.New("name, email, and phone are required fields")
	}
	return nil
}

func sendPasswordToEmail(user *domain.User) error {
	// Generate random password and hash
	plainPassword, hashedPassword, err := utils.GenerateRandomPasswordHash()
	if err != nil {
		return err
	}
	// Assign the hashed password to the user
	user.Password = hashedPassword

	// Send the plain password to the user's email
	emailSubject := "Welcome to UNO Service"
	emailBody := fmt.Sprintf("Hello %s,\n\nYour account has been created. Your temporary password is: %s\n\nPlease change it after logging in.", user.Name, plainPassword)

	err = utils.SendEmail(user.Email, emailSubject, emailBody)
	if err != nil {
		return errors.New("user created but failed to send email")
	}
	return nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (string, error) {

	if err := validateRequiredFields(user); err != nil {
		return "", err
	}
	// Generate UUID
	user.ID = uuid.New().String() // Generate UUID as string

	// Call the function to generate password, hash it, and send the email
	if err := sendPasswordToEmail(user); err != nil {
		return "", err
	}

	err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	return user.ID, nil
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID == "" {
		return errors.New("user ID is required")
	}

	if err := validateRequiredFields(user); err != nil {
		fmt.Println("Name, phone and e-mail cannot be empty.Previous values remained.")
	}

	if user.Password == "" {
		sendPasswordToEmail(user)
	}
	// If picture is empty in the body, set it to an empty string
	if user.Picture == "" {
		user.Picture = ""
	}

	err := s.Repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", user.ID, err)
	}

	return nil
}
