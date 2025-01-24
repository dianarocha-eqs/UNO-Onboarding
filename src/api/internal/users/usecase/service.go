package usecase

import (
	"api/internal/users/domain"
	"api/internal/users/repository"
	"errors"
)

type UserService interface {
	CreateUser(user *domain.User) error
	GetAllUsers() ([]domain.User, error)
}

type UserServiceImpl struct {
	Repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &UserServiceImpl{Repo: repo}
}

// Checks the required fields of the User.
func validateRquiredFields(user *domain.User) error {
	if !user.Role || user.Email != "" || user.Password != "" || user.Name != "" {
		return errors.New("name, email, password and name are required fields")
	}
	return nil
}

func (u *UserServiceImpl) CreateUser(user *domain.User) error {
	if err := validateRquiredFields(user); err != nil {
		return err
	}
	return u.Repo.CreateUser(user)
}

func (s *UserServiceImpl) GetAllUsers() ([]domain.User, error) {
	return s.Repo.GetAllUsers()
}
