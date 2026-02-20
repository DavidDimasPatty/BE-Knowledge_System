package repository

import (
	dto "be-knowledge/internal/delivery/dto/userManagement"
	"be-knowledge/internal/entities"
	Tracelog "be-knowledge/internal/tracelog"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/gomail.v2"
)

type UserManagementRepository interface {
	GetAllUser() (data *dto.UserManagement_GetAllUser_Response, er error)
	AddUser(data dto.UserManagement_AddUser_Request) error
	EditUserGet(id int) (data *entities.User, er error)
	EditUser(data dto.UserManagement_EditUser_Request) error
	DeleteUser(id int) error
	ChangeStatusUser(data dto.UserManagement_ChangeStatusUser_Request) error
}

type userManagementRepository struct {
	db *sqlx.DB
}

func NewUserManagementRepository(db *sqlx.DB) UserManagementRepository {
	return &userManagementRepository{db}
}

func (r *userManagementRepository) GetAllUser() (*dto.UserManagement_GetAllUser_Response, error) {
	namaEndpoint := "GetAllUser"
	res := dto.UserManagement_GetAllUser_Response{}

	users := []entities.User{}

	query := "SELECT u.*, r.nama AS `role_name` FROM users u LEFT JOIN roles r ON r.id = u.roles;"

	err := r.db.Select(&users, query)

	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s ", query),
		namaEndpoint,
	)
	if err != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("Error : %v", err.Error()),
			namaEndpoint,
		)
		return nil, err
	}

	res.Data = users

	return &res, nil
}

func (r *userManagementRepository) AddUser(data dto.UserManagement_AddUser_Request) error {
	namaEndpoint := "AddUser"
	var checkUser int
	querySelect := "SELECT count(*) FROM users WHERE username = ? or email = ?"

	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: username=%v, email=%v", querySelect, data.Username, data.Email),
		namaEndpoint,
	)
	errSelect := r.db.Get(&checkUser, querySelect, data.Username, data.Email)
	if errSelect != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("Error : %v", errSelect.Error()),
			namaEndpoint,
		)
		return errSelect
	}

	if checkUser > 0 {
		Tracelog.UserManagementLog(
			"Error : user already exist",
			namaEndpoint,
		)
		return errors.New("user already exist")
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	resultPassword := make([]byte, 8)
	for i := range resultPassword {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		resultPassword[i] = charset[num.Int64()]
	}

	passwordStr := string(resultPassword)

	hashedPass, errHash := bcrypt.GenerateFromPassword([]byte(passwordStr), bcrypt.DefaultCost)
	if errHash != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("Error : %v", errHash.Error()),
			namaEndpoint,
		)

		return errHash
	}

	query := `
    INSERT INTO users 
    (username, PASSWORD, email, nama, roles, passwordExpired, addId, divisi, STATUS, addTime)
    VALUES (?, ?, ?, ?, ?, DATE_ADD(NOW(), INTERVAL 30 YEAR), ?, 1, "Inactive", NOW())
`

	_, err := r.db.Exec(query,
		data.Username,
		string(hashedPass),
		data.Email,
		data.Nama,
		data.RoleId,
		data.AddId,
	)

	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: username=%v, email=%v, nama=%v, roleId=%v, addId=%v",
			err, data.Username, data.Email, data.Nama, data.RoleId, data.AddId),
		namaEndpoint)

	go sendEmail(data.Email, data.Username, passwordStr)

	return err
}

func (r *userManagementRepository) EditUserGet(id int) (*entities.User, error) {
	namaEndpoint := "EditUserGet"
	user := entities.User{}
	query := "SELECT * FROM users WHERE id = ?"

	err := r.db.Get(&user, query, id)
	if err != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
			namaEndpoint)
		return nil, err
	}

	return &user, nil
}

func (r *userManagementRepository) EditUser(data dto.UserManagement_EditUser_Request) error {
	namaEndpoint := "EditUser"
	var checkUser int
	querySelect := "SELECT count(*) FROM users WHERE ( email = ? ) and  id != ?"
	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: email=%v, id=%v", querySelect, data.Email, data.Id),
		namaEndpoint)

	errSelect := r.db.Get(&checkUser, querySelect, data.Email, data.Id)
	if errSelect != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("Error : %v", errSelect.Error()),
			namaEndpoint,
		)
		return errSelect
	}

	if checkUser > 0 {
		Tracelog.UserManagementLog(
			"Error : user already exist",
			namaEndpoint,
		)
		return errors.New("user already exist")
	}

	query := `
		UPDATE users 
		SET   email = ?, nama = ?, roles = ?, updId = ?, updTime = NOW()
		WHERE id = ?
	`
	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: email=%v, nama=%v, roles=%v, updId=%v, id=%v",
			query, data.Email, data.Nama, data.RoleId, data.UpdId, data.Id),
		namaEndpoint)

	_, err := r.db.Exec(query,
		data.Email,
		data.Nama,
		data.RoleId,
		data.UpdId,
		data.Id,
	)

	return err
}

func (r *userManagementRepository) DeleteUser(id int) error {
	namaEndpoint := "DeleteUser"
	query := "DELETE FROM users WHERE id = ?"
	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", query, id),
		namaEndpoint)
	_, err := r.db.Exec(query, id)
	return err
}

func (r *userManagementRepository) ChangeStatusUser(data dto.UserManagement_ChangeStatusUser_Request) error {
	namaEndpoint := "ChangeStatusUser"
	var checkUser int
	querySelect := "SELECT count(*) FROM users WHERE  id = ?"

	errSelect := r.db.Get(&checkUser, querySelect, data.Id)

	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: id=%v", querySelect, data.Id),
		namaEndpoint)

	if errSelect != nil {
		Tracelog.UserManagementLog(
			fmt.Sprintf("Error : %v", errSelect.Error()),
			namaEndpoint,
		)
		return errSelect
	}

	if checkUser == 0 {
		Tracelog.UserManagementLog(
			"Error : user not found",
			namaEndpoint,
		)
		return errors.New("user not found")
	}

	query := "UPDATE users set status = ? where id = ?"

	var statusFinal string
	if data.Status == "Active" {
		statusFinal = "Block"
	} else {
		statusFinal = "Active"
	}
	Tracelog.UserManagementLog(
		fmt.Sprintf("SQL: %s | Params: status=%v, id=%v", querySelect, data.Id, statusFinal),
		namaEndpoint)
	_, err := r.db.Exec(query, statusFinal, data.Id)
	return err
}

func sendEmail(to, username, password string) error {
	m := gomail.NewMessage()

	m.SetHeader("From", os.Getenv("SMTP_FROM"))
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Your Account Information")

	body := fmt.Sprintf(
		"Your account has been created.\n\nUsername: %s\nPassword: %s\n\nPlease change your password after login.",
		username,
		password,
	)

	m.SetBody("text/plain", body)
	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		return err
	}

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"),
		port,
		os.Getenv("SMTP_USER"),
		os.Getenv("SMTP_PASS"))

	return d.DialAndSend(m)
}
