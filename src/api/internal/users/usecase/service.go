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

	TrimUserFields(user)

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

// trim spaces for required fields only
// reflection is more slow than Join
func TrimUserFields(user *domain.User) {
	if user == nil {
		return
	}

	user.Name = strings.Join(strings.Fields(user.Name), " ")
	user.Email = strings.Join(strings.Fields(user.Email), " ")
	user.Phone = strings.Join(strings.Fields(user.Phone), " ")
}

// Creates the password and sends it to the user via email
func createPassword(user *domain.User, password string) (string, error) {
	// Generate password hash (either random or user-provided)
	var plainPassword string
	var hashedPassword string
	var err error

	plainPassword, hashedPassword, err = utils.GeneratePasswordHash(password)
	if err != nil {
		return "", err
	}
	// Assign the hashed password to the user
	user.Password = hashedPassword

	return plainPassword, nil
}

func sendEmail(user *domain.User, plainPassword string) error {
	// Send the plain password to the user's email
	var emailSubject string
	var emailBody string
	emailSubject = "Welcome to UNO Service"
	emailBody = fmt.Sprintf("Hello %s,\n\nYour account has been created. Your temporary password is: %s\n\nPlease change it after logging in.", user.Name, plainPassword)

	err := utils.SendEmail(user.Email, emailSubject, emailBody)
	if err != nil {
		return errors.New("user created but failed to send email")
	}
	return nil
}

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (uuid.UUID, error) {
	var err error
	if err = validateRequiredFields(user); err != nil {
		return uuid.NilUUID, err
	}

	// Create the password (and hash it)
	var plainPasswordForEmail string
	plainPasswordForEmail, err = createPassword(user, "")
	if err != nil {
		return uuid.NilUUID, err
	}

	err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return uuid.NilUUID, errors.New("failed to create user")
	}
	fmt.Println(plainPasswordForEmail)

	// Send the email with the plain password (only after user is created)
	if err = sendEmail(user, plainPasswordForEmail); err != nil {
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

func (s *UserServiceImpl) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	if userID == (uuid.UUID{}) {
		return nil, errors.New("user ID is required")
	}

	// Fetch user from repository
	user, err := s.Repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}

func (s *UserServiceImpl) GetUserByEmailAndPassword(ctx context.Context, email, password string) (*domain.User, error) {
	if email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if password == "" {
		return nil, fmt.Errorf("password cannot be empty")
	}

	user, err := s.Repo.GetUserByEmailAndPassword(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user: %v", err)
	}

	return user, nil
}
