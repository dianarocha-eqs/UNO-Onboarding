package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"api/util"
	"context"
	"errors"
	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *domain.User) (string, error)
	// GetAllUsers() ([]domain.User, error)
}

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

// Checks the required fields of the User.
func validateRquiredFields(user *domain.User) error {
	if user.Email == "" || user.Name == "" || user.Phone == "" {
		return errors.New("name, email, and password are required fields")
	}
	return nil
}

// func (s *UserServiceImpl) ValidateAdmin(ctx context.Context) error {
// 	admin, err := s.Repo.FindUserByRole(ctx, true) // Find admin
// 	if err != nil || admin == nil {
// 		return errors.New("403: only admins can perform this action")
// 	}
// 	return nil
// }

func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) (string, error) {
	// Validate admin privileges
	// err := s.ValidateAdmin(ctx)
	// if err != nil {
	// 	return "", err
	// }

	if err := validateRquiredFields(user); err != nil {
		return "", err
	}
	// Generate random password and hash
	plainPassword, hashedPassword, err := util.GenerateRandomPasswordHash()
	if err != nil {
		return "", err
	}

	// Assign the hashed password to the user model
	user.Password = hashedPassword

	// Save the user
	err = s.Repo.CreateUser(ctx, user)
	if err != nil {
		return "", err
	}
	// Send plain password via email (future implementation)
	// For now, log the plain password (temporary)
	fmt.Printf("User created successfully. Password: %s\n", plainPassword)

	return user.ID, nil
}

// func (s *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
// 	return s.Repo.GetAllUsers()
// }
