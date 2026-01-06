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
	EditPassword(username string, newPassword string, oldPassword string) error
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

	if user.Status != nil && *user.Status == "Inactive" {
		return nil, errors.New("akun anda belum aktif, hubungi admin")
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

func (s *userService) EditPassword(username string, newPassword string, oldPassword string) error {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return errors.New("username tidak ditemukan")
	}

	// Cek apakah user sudah diblock
	if user.Status != nil && *user.Status == "Block" {
		return errors.New("akun anda diblok, hubungi admin")
	}

	if user.Status != nil && *user.Status == "Inactive" {
		return errors.New("akun anda belum aktif, hubungi admin")
	}

	// Cek password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		return errors.New("password tidak boleh sama dengan password sekarang")
	}

	if user.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*user.OldPassword), []byte(newPassword))
		if err == nil {
			return errors.New("password tidak boleh sama dengan password sebelumnya")
		}
	}

	// --- Password benar ---
	err = s.repo.ChangePassword(username, newPassword, oldPassword)
	if err != nil {
		return errors.New("Internal Error")
	}

	return nil
}
