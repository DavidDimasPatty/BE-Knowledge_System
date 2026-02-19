package usecases

import (
	"be-knowledge/configs"
	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
	Tracelog "be-knowledge/internal/tracelog"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type UserService interface {
	Login(username, password string) (*entities.User, error)
	EditPassword(username string, newPassword string, oldPassword string) error
	SendEmailResetPassword(email string) error
	ValidateResetToken(token string) error
	ResetPassword(token, newPassword string) error
	AuthUser(userID int) error
}

type userService struct {
	repo         repository.UserRepository
	config       *configs.Config
	emailService EmailService
}

func NewUserService(repo repository.UserRepository, cfg *configs.Config, emailService EmailService) UserService {
	return &userService{
		repo:         repo,
		config:       cfg,
		emailService: emailService,
	}
}

func (s *userService) Login(username, password string) (*entities.User, error) {
	now := time.Now()

	var namaEndpoint = "Login"

	Tracelog.AuthLog(fmt.Sprintf("Get Username : %v", username), namaEndpoint)

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		Tracelog.AuthLog(fmt.Sprintf("Gagal username tidak ditemukan : %v", username), namaEndpoint)
		return nil, errors.New("username tidak ditemukan")
	}

	// Cek apakah user sudah diblock
	if user.Status != nil && *user.Status == "Block" {
		Tracelog.AuthLog(fmt.Sprintf("Gagal akun anda diblok, hubungi admin : %v", username), namaEndpoint)
		return nil, errors.New("akun anda diblok, hubungi admin")
	}

	if user.Status != nil && *user.Status == "Inactive" {
		Tracelog.AuthLog(fmt.Sprintf("Gagal akun anda belum aktif, hubungi admin : %v", username), namaEndpoint)
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
			Tracelog.AuthLog(fmt.Sprintf("Gagal akun anda diblok karena terlalu banyak salah password : %v", username), namaEndpoint)
			return nil, errors.New("akun anda diblok karena terlalu banyak salah password")
		}

		return nil, errors.New("password salah")
	}

	// --- Password benar ---

	// Reset loginCount
	_ = s.repo.ResetLoginCount(user.ID)

	err = s.repo.UpdateLastLogin(user.ID, now)
	if err != nil {
		Tracelog.AuthLog(fmt.Sprintf("Gagal update last login : %v", username), namaEndpoint)
		return nil, errors.New("gagal update last login")
	}

	user.LastLogin = &now
	return user, nil
}

func (s *userService) EditPassword(username string, newPassword string, oldPassword string) error {
	namaEndpoint := "EditPassword"
	Tracelog.AuthLog("Mulai proses EditPassword", namaEndpoint)
	Tracelog.AuthLog(fmt.Sprintf("Parameter -> username: %s", username), namaEndpoint)

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		Tracelog.AuthLog("Username tidak ditemukan: "+err.Error(), namaEndpoint)
		return errors.New("username tidak ditemukan")
	}

	// Cek apakah user sudah diblock
	if user.Status != nil && *user.Status == "Block" {
		Tracelog.AuthLog("Akun diblok", namaEndpoint)
		return errors.New("akun anda diblok, hubungi admin")
	}

	if user.Status != nil && *user.Status == "Inactive" {
		Tracelog.AuthLog("Akun belum aktif", namaEndpoint)
		return errors.New("akun anda belum aktif, hubungi admin")
	}

	// Cek password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if err != nil {
		Tracelog.AuthLog("Password lama salah", namaEndpoint)
		return errors.New("password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		Tracelog.AuthLog("Password baru sama dengan password sekarang", namaEndpoint)
		return errors.New("password tidak boleh sama dengan password sekarang")
	}

	if user.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*user.OldPassword), []byte(newPassword))
		if err == nil {
			Tracelog.AuthLog("Password baru sama dengan password sebelumnya", namaEndpoint)
			return errors.New("password tidak boleh sama dengan password sebelumnya")
		}
	}

	// --- Password benar ---
	err = s.repo.ChangePassword(username, newPassword, oldPassword)
	if err != nil {
		Tracelog.AuthLog("Gagal update password: "+err.Error(), namaEndpoint)
		return errors.New("Internal Error")
	}

	Tracelog.AuthLog("Berhasil update password", namaEndpoint)
	Tracelog.AuthLog("Selesai proses EditPassword", namaEndpoint)

	return nil
}

func (s *userService) SendEmailResetPassword(email string) error {
	namaEndpoint := "SendEmailResetPassword"
	Tracelog.AuthLog("Mulai proses SendEmailResetPassword", namaEndpoint)
	Tracelog.AuthLog(fmt.Sprintf("Parameter -> email: %s", email), namaEndpoint)

	user, err := s.repo.GetByEmail(email)
	if err != nil {
		Tracelog.AuthLog("Email tidak ditemukan", namaEndpoint)
		return errors.New("")
		// return nil
	}

	// Cek apakah user sudah diblock
	if user.Status != nil && *user.Status == "Block" {
		Tracelog.AuthLog("Akun diblok", namaEndpoint)
		return errors.New("akun anda diblok, hubungi admin")
	}

	if user.Status != nil && *user.Status == "Inactive" {
		Tracelog.AuthLog("Akun belum aktif", namaEndpoint)
		return errors.New("akun anda belum aktif, hubungi admin")
	}

	token, err := generateResetToken()
	if err != nil {
		Tracelog.AuthLog("Gagal generate token: "+err.Error(), namaEndpoint)
		return err
	}

	addTime := time.Now()
	expiredDate := addTime.Add(30 * time.Minute)

	err = s.repo.SaveResetToken(user.ID, token, expiredDate, addTime)
	if err != nil {
		Tracelog.AuthLog("Gagal simpan token reset: "+err.Error(), namaEndpoint)
		return err
	}

	resetLink := fmt.Sprintf(
		"%s/reset-password?token=%s",
		s.config.FrontendURL,
		token,
	)

	subject := "Reset Password"
	body := fmt.Sprintf(
		"Halo %s,\n\n"+
			"Klik link berikut untuk reset password:\n\n%s\n\n"+
			"Link ini berlaku selama 30 menit.\n\n"+
			"Jika Anda tidak merasa melakukan permintaan ini, abaikan email ini.",
		user.Username,
		resetLink,
	)

	if user.Email == nil || *user.Email == "" {
		Tracelog.AuthLog("Email user tidak valid", namaEndpoint)
		return errors.New("email user tidak valid")
	}

	return s.emailService.SendEmail(
		*user.Email,
		subject,
		body,
	)
}

func generateResetToken() (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (s *userService) ValidateResetToken(token string) error {
	namaEndpoint := "ValidateResetToken"
	Tracelog.AuthLog("Mulai proses ValidateResetToken", namaEndpoint)

	data, err := s.repo.GetResetToken(token)
	if err != nil {
		Tracelog.AuthLog("Token tidak valid", namaEndpoint)
		return errors.New("token tidak valid")
	}

	if *data.IsReset == "Y" {
		Tracelog.AuthLog("Token sudah digunakan", namaEndpoint)
		return errors.New("token sudah digunakan")
	}

	if time.Now().After(*data.ExpiredDate) {
		Tracelog.AuthLog("Token sudah kadaluarsa", namaEndpoint)
		return errors.New("token sudah kadaluarsa")
	}

	Tracelog.AuthLog("Token valid", namaEndpoint)
	Tracelog.AuthLog("Selesai proses ValidateResetToken", namaEndpoint)

	return nil
}

func (s *userService) ResetPassword(token string, newPassword string) error {
	namaEndpoint := "ResetPassword"
	Tracelog.AuthLog("Mulai proses ResetPassword", namaEndpoint)

	// Ambil token
	resetData, err := s.repo.GetResetToken(token)
	if err != nil {
		Tracelog.AuthLog("Token tidak valid", namaEndpoint)
		return errors.New("token tidak valid")
	}

	// Validasi token
	if *resetData.IsReset == "Y" {
		Tracelog.AuthLog("Token sudah digunakan", namaEndpoint)
		return errors.New("token sudah digunakan")
	}

	if time.Now().After(*resetData.ExpiredDate) {
		Tracelog.AuthLog("Token kadaluarsa", namaEndpoint)
		return errors.New("token sudah kedaluwarsa")
	}

	// Ambil user
	user, err := s.repo.GetById(resetData.User_Id)
	if err != nil {
		Tracelog.AuthLog("User tidak ditemukan", namaEndpoint)
		return errors.New("user tidak ditemukan")
	}

	// cegah password sama dengan password sekarang
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(newPassword))
	if err == nil {
		Tracelog.AuthLog("Password baru tidak boleh sama dengan password sekarang", namaEndpoint)
		return errors.New("password baru tidak boleh sama dengan password sekarang")
	}

	// cegah password sama dengan password sebelumnya
	if user.OldPassword != nil {
		err = bcrypt.CompareHashAndPassword([]byte(*user.OldPassword), []byte(newPassword))
		if err == nil {
			Tracelog.AuthLog("Password baru tidak boleh sama dengan password sebelumnya", namaEndpoint)
			return errors.New("password baru tidak boleh sama dengan password sebelumnya")
		}
	}

	// Ganti password
	err = s.repo.ChangePasswordByReset(
		user.Username,
		newPassword,
		user.Password, // dummy old password
	)
	if err != nil {
		Tracelog.AuthLog("Gagal reset password: "+err.Error(), namaEndpoint)
		return errors.New("gagal reset password")
	}

	// Tandai token sudah digunakan
	err = s.repo.MarkResetTokenUsed(token)
	if err != nil {
		Tracelog.AuthLog("Gagal update token: "+err.Error(), namaEndpoint)
		return errors.New("gagal update token")
	}

	Tracelog.AuthLog("Berhasil reset password", namaEndpoint)
	Tracelog.AuthLog("Selesai proses ResetPassword", namaEndpoint)

	return nil
}

func (s *userService) AuthUser(userID int) error {
	namaEndpoint := "AuthUser"
	Tracelog.AuthLog("Mulai proses AuthUser", namaEndpoint)
	Tracelog.AuthLog(fmt.Sprintf("Parameter -> userID: %d", userID), namaEndpoint)

	err := s.repo.GetByUserId(userID)
	if err != nil {
		Tracelog.AuthLog("User tidak ditemukan: "+err.Error(), namaEndpoint)
		return errors.New("username tidak ditemukan")
	}

	Tracelog.AuthLog("User ditemukan", namaEndpoint)
	Tracelog.AuthLog("Selesai proses AuthUser", namaEndpoint)

	return nil
}
