package repository

import (
	dto "be-knowledge/internal/delivery/dto/userManagement"
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"
)

type UserManagementRepository interface {
	GetAllUser() (data *dto.UserManagement_GetAllUser_Response, er error)
	AddUser(data *dto.UserManagement_AddUser_Request) error
	EditUserGet(id int) (data *entities.User, er error)
	EditUser(data *dto.UserManagement_EditUser_Request) error
	DeleteUser(id int) error
}

type userManagementRepository struct {
	db *sqlx.DB
}

func NewUserManagementRepository(db *sqlx.DB) UserManagementRepository {
	return &userManagementRepository{db}
}

func (r *userManagementRepository) GetAllUser() (*dto.UserManagement_GetAllUser_Response, error) {
	res := dto.UserManagement_GetAllUser_Response{}

	users := []entities.User{}

	query := "SELECT * FROM users"

	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	res.Data = users

	return &res, nil
}

func (r *userManagementRepository) AddUser(data *dto.UserManagement_AddUser_Request) error {
	query := `
		INSERT INTO users 
		(email, noTelp, nama, roles, addId, STATUS, addTime)
		VALUES (?, ?, ?, ?, ?, "Inactive" ,NOW())
	`

	_, err := r.db.Exec(query,
		data.Email,
		data.NoTelp,
		data.Nama,
		data.RoleId,
		data.AddId,
	)

	return err
}

func (r *userManagementRepository) EditUserGet(id int) (*entities.User, error) {
	user := entities.User{}
	query := "SELECT * FROM users WHERE id = ?"

	err := r.db.Get(&user, query, id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userManagementRepository) EditUser(data *dto.UserManagement_EditUser_Request) error {
	query := `
		UPDATE users 
		SET  email = ?, noTelp = ?, nama = ?, roles = ?, updId = ?, updTime = NOW()
		WHERE id = ?
	`

	_, err := r.db.Exec(query,
		data.Email,
		data.NoTelp,
		data.Nama,
		data.RoleId,
		data.UpdId,
		data.Id,
	)

	return err
}

func (r *userManagementRepository) DeleteUser(id int) error {
	query := "DELETE FROM users WHERE id = ?"

	_, err := r.db.Exec(query, id)
	return err
}
