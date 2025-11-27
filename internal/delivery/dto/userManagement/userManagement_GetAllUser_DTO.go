package dto

import "be-knowledge/internal/entities"

type UserManagement_GetAllUser_Response struct {
	Data []entities.User `json:"data"`
}

type UserManagement_GetAllUser_Request struct {
	Id     int    `json:"id"`
	Nama   string `json:"nama"`
	RoleId int    `json:"roleId"`
	Email  string `json:"email"`
	NoTelp string `json:"noTelp"`
}
