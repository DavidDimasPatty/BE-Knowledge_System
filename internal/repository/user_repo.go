package repository

import (
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"

	"time"
)

type UserRepository interface {
	GetByUsername(username string) (*entities.User, error)
	UpdateLastLogin(id int, lastLogin time.Time) error
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetByUsername(username string) (*entities.User, error) {
	user := entities.User{}
	query := "SELECT * FROM users WHERE username = ? LIMIT 1"

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
