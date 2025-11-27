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
	user.LastLogin = &now

	return user, nil
}
