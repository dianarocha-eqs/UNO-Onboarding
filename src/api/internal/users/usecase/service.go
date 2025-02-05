package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/utils"
	"context"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"

	uuid "github.com/tentone/mssql-uuid"
)

// Interface for user's services
type UserService interface {
	// Creates a new user and returns the user's UUID
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	// Updates an existing user
	UpdateUser(ctx context.Context, user *domain.User) error
	// Get users
	GetUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error)
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
	// trim blank spaces for each field
	TrimUserFields(user)

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

// Trims extra spaces from all string fields of a given user
func TrimUserFields(user *domain.User) {
	if user == nil {
		return
	}

	v := reflect.ValueOf(user).Elem()

	// Loop through all fields
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)

		// Check if the field is a string
		if field.Kind() == reflect.String {
			field.SetString(strings.Join(strings.Fields(field.String()), " "))
		}
	}
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
	user.ID = uuid.NewV4() // Generate UUID

	if err := sendPasswordToEmail(user, ""); err != nil {
		return "", err
	}

	err := s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", errors.New("failed to create user")
	}
	return user.ID.String(), nil // Convert UUID to string before returning
}

func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	if user.ID == (uuid.UUID{}) {
		return errors.New("user ID is required")
	}

	// Get current user data from the database using UUID
	currentUser, err := s.Repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return fmt.Errorf("failed to retrieve current user: %v", err)
	}

	// trim blank spaces for each field
	TrimUserFields(user)

	// Preserve existing values if not provided
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

	// Hash the password if provided, otherwise keep the current one
	if user.Password != "" {
		_, hashedPassword, err := utils.GeneratePasswordHash(user.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPassword
	} else {
		user.Password = currentUser.Password
	}

	// Update the user in the database
	err = s.Repo.UpdateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user with id %s: %v", user.ID.String(), err)
	}
	return nil
}

func (s *UserServiceImpl) GetUsers(ctx context.Context, search string, sortDirection int) ([]domain.User, error) {

	var users []domain.User
	var err error
	users, err = s.Repo.GetUsers(ctx)
	if err != nil {
		return nil, errors.New("failed to retrieve users")
	}

	// fmt.Println(users)

	search = strings.Join(strings.Fields(search), " ")
	// Filter by name or email
	if search != "" {
		var filteredUsers []domain.User
		for _, user := range users {
			if strings.Contains(user.Name, search) || strings.Contains(user.Email, search) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		// if there is no search on name and email that matches, then filteredUsers has nothing
		if len(filteredUsers) == 0 {
			return nil, errors.New("no matching user found: the search term does not appear in either name or email")
		}

		users = filteredUsers
	}

	if sortDirection != 1 && sortDirection != -1 && sortDirection != 0 {
		return nil, errors.New("sort direction value is wrong ( different from 1, -1 or 0)")
	}

	// Sort by name
	if sortDirection == 1 {
		// Ascending
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name < users[j].Name
		})
	} else if sortDirection == -1 {
		// Descending
		sort.Slice(users, func(i, j int) bool {
			return users[i].Name > users[j].Name
		})
	}

	return users, nil
}
