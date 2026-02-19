package repository

import (
	"be-knowledge/internal/entities"
	Tracelog "be-knowledge/internal/tracelog"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByUsername(username string) (*entities.User, error)
	GetByUserId(id int) error
	UpdateLastLogin(id int, lastLogin time.Time) error
	BlockUser(id int, blockDate time.Time) error
	IncrementLoginCount(id int) error
	ResetLoginCount(id int) error
	ChangePassword(username string, newPassword string, oldPassword string) error
	GetByEmail(email string) (*entities.User, error)
	SaveResetToken(userID int, token string, expiredDate time.Time, addTime time.Time) error
	GetResetToken(token string) (*entities.PasswordResets, error)
	MarkResetTokenUsed(token string) error
	GetById(id int) (*entities.User, error)
	ChangePasswordByReset(username string, newPassword string, currentHashedPassword string) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	var namaEndpoint = "GetByUsername"
	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE username = ? LIMIT 1`

	err := r.db.Get(&user, query, username)
	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: username=%v", query, username),
		namaEndpoint,
	)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetByUserId(id int) error {
	var namaEndpoint = "GetByUserId"
	var count int
	query := `
		SELECT count(*)
		FROM users
		WHERE id = ? LIMIT 1`

	err := r.db.Get(&count, query, id)
	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint,
	)
	if err != nil {
		return err
	}

	if count == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (r *userRepository) UpdateLastLogin(id int, lastLogin time.Time) error {
	var namaEndpoint = "UpdateLastLogin"

	query := "UPDATE users SET lastLogin = ? WHERE id = ?"

	_, err := r.db.Exec(query, lastLogin, id)
	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: lastLogin=%v, id=%v", query, lastLogin, id),
		namaEndpoint,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) BlockUser(id int, blockDate time.Time) error {
	var namaEndpoint = "BlockUser"

	query := "UPDATE users SET status = 'Block', blockDate = ? WHERE id = ?"

	_, err := r.db.Exec(query, blockDate, id)
	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: blockDate=%v, id=%v", query, blockDate, id),
		namaEndpoint,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) IncrementLoginCount(id int) error {
	var namaEndpoint = "IncrementLoginCount"

	query := "UPDATE users SET loginCount = loginCount + 1 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint,
	)
	return err
}

func (r *userRepository) ResetLoginCount(id int) error {
	var namaEndpoint = "ResetLoginCount"

	query := "UPDATE users SET loginCount = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint,
	)
	return err
}

func (r *userRepository) ChangePassword(username string, newPassword string, oldPassword string) error {
	var namaEndpoint = "ChangePassword"

	hashedPass, errHash := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return errHash
	}

	hashedOldPass, errHash := bcrypt.GenerateFromPassword([]byte(oldPassword), bcrypt.DefaultCost)
	if errHash != nil {
		return errHash
	}

	query := `
    update users set password=?, oldPassword=?, passwordExpired=DATE_ADD(NOW(), INTERVAL 30 YEAR)
	where username=?`

	_, err := r.db.Exec(query,
		string(hashedPass),
		string(hashedOldPass),
		username,
	)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: username=%v", query, username),
		namaEndpoint,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	var namaEndpoint = "GetByEmail"

	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE email = ? LIMIT 1`

	err := r.db.Get(&user, query, email)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: email=%v", query, email),
		namaEndpoint,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveResetToken(userID int, token string, expiredDate time.Time, addTime time.Time) error {
	var namaEndpoint = "SaveResetToken"

	query := `
		INSERT INTO passwordresets (user_id, token, expiredDate, addTime, isReset)
		VALUES (?, ?, ?, ?, 'N')
	`
	_, err := r.db.Exec(query, userID, token, expiredDate, addTime)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: userID=%v, token=%v, expiredDate=%v, addTime=%v",
			query, userID, token, expiredDate, addTime),
		namaEndpoint,
	)

	return err
}

func (r *userRepository) GetResetToken(token string) (*entities.PasswordResets, error) {
	var namaEndpoint = "GetResetToken"

	data := entities.PasswordResets{}

	query := `
		SELECT *
		FROM passwordresets
		WHERE token = ?
		LIMIT 1
	`

	err := r.db.Get(&data, query, token)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: token=%v", query, token),
		namaEndpoint,
	)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *userRepository) MarkResetTokenUsed(token string) error {
	var namaEndpoint = "MarkResetTokenUsed"

	query := `
		UPDATE passwordresets
		SET isReset = 'Y'
		WHERE token = ?
	`

	_, err := r.db.Exec(query, token)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: token=%v", query, token),
		namaEndpoint,
	)

	return err
}

func (r *userRepository) GetById(id int) (*entities.User, error) {
	var namaEndpoint = "GetById"

	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE u.id = ? LIMIT 1`

	err := r.db.Get(&user, query, id)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ChangePasswordByReset(username string, newPassword string, currentHashedPassword string) error {
	var namaEndpoint = "ChangePasswordByReset"

	hashedNew, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	query := `
		UPDATE users
		SET 
			oldPassword = ?,
			password = ?,
			passwordExpired = DATE_ADD(NOW(), INTERVAL 30 YEAR)
		WHERE username = ?
	`

	_, err = r.db.Exec(
		query,
		currentHashedPassword,
		string(hashedNew),
		username,
	)

	Tracelog.AuthLog(
		fmt.Sprintf("SQL: %s | Params: username=%v", query, username),
		namaEndpoint,
	)

	return err
}
