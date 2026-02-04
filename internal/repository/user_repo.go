package repository

import (
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"

	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetByUsername(username string) (*entities.User, error)
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
	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE username = ? LIMIT 1`

	err := r.db.Get(&user, query, username)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) UpdateLastLogin(id int, lastLogin time.Time) error {
	query := "UPDATE users SET lastLogin = ? WHERE id = ?"

	_, err := r.db.Exec(query, lastLogin, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) BlockUser(id int, blockDate time.Time) error {
	query := "UPDATE users SET status = 'Block', blockDate = ? WHERE id = ?"

	_, err := r.db.Exec(query, blockDate, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) IncrementLoginCount(id int) error {
	query := "UPDATE users SET loginCount = loginCount + 1 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) ResetLoginCount(id int) error {
	query := "UPDATE users SET loginCount = 0 WHERE id = ?"
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userRepository) ChangePassword(username string, newPassword string, oldPassword string) error {
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

	if err != nil {
		return err
	}

	return nil
}

func (r *userRepository) GetByEmail(email string) (*entities.User, error) {
	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE email = ? LIMIT 1`

	err := r.db.Get(&user, query, email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) SaveResetToken(userID int, token string, expiredDate time.Time, addTime time.Time) error {
	query := `
		INSERT INTO passwordresets (user_id, token, expiredDate, addTime, isReset)
		VALUES (?, ?, ?, ?, 'N')
	`
	_, err := r.db.Exec(query, userID, token, expiredDate, addTime)
	return err
}

func (r *userRepository) GetResetToken(token string) (*entities.PasswordResets, error) {
	data := entities.PasswordResets{}

	query := `
		SELECT *
		FROM passwordresets
		WHERE token = ?
		LIMIT 1
	`

	err := r.db.Get(&data, query, token)
	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *userRepository) MarkResetTokenUsed(token string) error {
	query := `
		UPDATE passwordresets
		SET isReset = 'Y'
		WHERE token = ?
	`

	_, err := r.db.Exec(query, token)
	return err
}

func (r *userRepository) GetById(id int) (*entities.User, error) {
	user := entities.User{}
	query := `
		SELECT u.*, r.nama AS role_name
		FROM users u 
		LEFT JOIN roles r ON u.roles = r.id
		WHERE u.id = ? LIMIT 1`

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) ChangePasswordByReset(username string, newPassword string, currentHashedPassword string) error {
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

	return err
}
