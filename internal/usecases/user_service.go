package usecases

import (
	"errors"
	"time"

	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(username, password string) (*entities.User, error)
	UpdateLastLogin(id int, t time.Time) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) Login(username, password string) (*entities.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, errors.New("username tidak ditemukan")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("password salah")
	}

	now := time.Now()

	// Update last login di DB
    err = s.UpdateLastLogin(user.ID, now)
    if err != nil {
        return nil, errors.New("gagal update last login")
    }

	user.LastLogin = &now

	return user, nil
}

func (s *userService) UpdateLastLogin(id int, t time.Time) error {
    return s.repo.UpdateLastLogin(id, t)
}