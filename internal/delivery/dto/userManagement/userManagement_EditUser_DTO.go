package dto

type UserManagement_EditUser_Request struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	RoleId int    `json:"roleId"`
	Email  string `json:"email"`
	UpdId  string `json:"updId"`
	// NoTelp string `json:"noTelp"`
}
