package dto

type UserManagement_AddUser_Request struct {
	Nama     string `json:"nama"`
	Username string `json:"username"`
	RoleId   int    `json:"roleId"`
	Email    string `json:"email"`
	NoTelp   string `json:"noTelp"`
	AddId    string `json:"addId"`
}
