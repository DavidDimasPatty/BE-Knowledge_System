package dto

import "be-knowledge/internal/entities"

type Auth_Login_Response struct {
	Data []entities.User `json:"data"`
}

type Auth_Login_Request struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
