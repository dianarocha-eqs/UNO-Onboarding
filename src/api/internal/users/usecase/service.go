package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

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

	// takes spaces
	user.Name = strings.Join(strings.Fields(user.Name), " ")
	user.Email = strings.Join(strings.Fields(user.Email), " ")
	user.Phone = strings.Join(strings.Fields(user.Phone), " ")

	if user.Email == "" || user.Name == "" || user.Phone == "" {
		return errors.New("name, email, and phone are required fields")
	}

	// Validate email
	if err := validateEmail(user.Email); err != nil {
		return err
	}
	// Validate phone number
	if err := validatePhone(user.Phone); err != nil {
		return err
	}

	return nil
}

// Checks if the email is well constructed
func validateEmail(email string) error {
	var checkemail = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !checkemail.MatchString(email) || len(email) < 5 {
		return errors.New("invalid email format")
	}
	return nil
}

// Checks if the phone is well constructed
func validatePhone(phone string) error {
	var checknumber = regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)
	if !checknumber.MatchString(phone) {
		return errors.New("invalid phone number format")
	}
	return nil
}

// Creates the password and sends it to the user via email
func sendPasswordToEmail(user *domain.User, password string) error {
	// Generate password hash (either random or user-provided)
	plainPassword, hashedPassword, err := utils.GeneratePasswordHash(password)
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
	// The password is empty, so it will generate a random password
	if err := sendPasswordToEmail(user, ""); err != nil {
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

	// Fetch the current user data from the database
	currentUser, err := s.Repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve current user: %v", err)
	}

	// If one of these three fields is empty, it remains as the previous values (should never be empty)
	if user.Name == "" {
		user.Name = currentUser.Name
	}
	if user.Email == "" {
		user.Email = currentUser.Email
	}
	if user.Phone == "" {
		user.Phone = currentUser.Phone
	}

	// Picture can be empty
	if user.Picture == "" {
		user.Picture = ""
	}

	// If the password is provided in the request body, hash and update it
	if user.Password != "" {
		// Hash the password provided by the user
		_, hashedPassword, err := utils.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	} else {
		// If password is not provided, retain the existing password
		user.Password = currentUser.Password
	}

	err = s.Repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", user.ID, err)
	}

	return nil
}
