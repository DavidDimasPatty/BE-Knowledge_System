package dto

type UserManagement_EditUser_Request struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	RoleId int    `json:"roleId"`
	Email  string `json:"email"`
	NoTelp string `json:"noTelp"`
	UpdId  string `json:"updId"`
}
