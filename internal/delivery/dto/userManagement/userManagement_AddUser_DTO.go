package dto

type UserManagement_AddUser_Request struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	RoleId int    `json:"roleId"`
	Email  string `json:"email"`
	NoTelp string `json:"noTelp"`
	AddId  string `json:"addId"`
}
