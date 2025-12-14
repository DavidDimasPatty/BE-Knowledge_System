package dto

import "be-knowledge/internal/entities"

type Category_GetAllCategoryUser_Response struct {
	Data []entities.Category `json:"data"`
}

type Category_GetAllCategoryUser_Request struct {
	Username string  `form:"username" binding:"required"`
	Search   *string `form:"search"`
	Page     *int    `form:"page"`
	Limit    *int    `form:"limit"`
}
