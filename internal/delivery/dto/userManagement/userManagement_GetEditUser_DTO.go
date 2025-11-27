package dto

import "be-knowledge/internal/entities"

type UserManagement_GetEditUser_Response struct {
	Data entities.User `json:"data"`
}

type UserManagement_GetEditUser_Request struct {
	Id int `json:"id"`
}
