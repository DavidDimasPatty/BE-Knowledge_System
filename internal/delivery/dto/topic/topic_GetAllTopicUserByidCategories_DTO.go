package dto

import "be-knowledge/internal/entities"

type Topic_GetAllTopicUserByidCategories_Response struct {
	Data []entities.Topic `json:"data"`
}

type Topic_GetAllTopicUserByidCategories_Request struct {
	Username     string `form:"username" binding:"required"`
	IdCategories int    `form:"idCategories" binding:"required"`
}
