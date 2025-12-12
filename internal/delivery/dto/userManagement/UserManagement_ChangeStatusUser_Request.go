package dto

type UserManagement_ChangeStatusUser_Request struct {
	Id     int    `json:"id"`
	Status string `json:"Status"`
}
