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
	now := time.Now()

    user, err := s.repo.GetByUsername(username)
    if err != nil {
        return nil, errors.New("username tidak ditemukan")
    }

    // Cek apakah user sudah diblock
    if user.Status != nil && *user.Status == "Block" {
        return nil, errors.New("akun anda diblok, hubungi admin")
    }

    // Cek password
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {

        // Tambah loginCount
        _ = s.repo.IncrementLoginCount(user.ID)

        // Ambil ulang user untuk mengecek loginCount terbaru
        user, _ = s.repo.GetByUsername(username)

        if user.LoginCount >= 5 {
            _ = s.repo.BlockUser(user.ID, now)
            return nil, errors.New("akun anda diblok karena terlalu banyak salah password")
        }

        return nil, errors.New("password salah")
    }

    // --- Password benar ---

    // Reset loginCount
    _ = s.repo.ResetLoginCount(user.ID)

    err = s.repo.UpdateLastLogin(user.ID, now)
    if err != nil {
        return nil, errors.New("gagal update last login")
    }

    user.LastLogin = &now
    return user, nil
}
