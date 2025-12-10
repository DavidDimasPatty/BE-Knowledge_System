package dto

import "be-knowledge/internal/entities"

type Topic_GetAllTopicUser_Response struct {
	Data []entities.Topic `json:"data"`
}

type Topic_GetAllTopicUser_Request struct {
	Username string `form:"username" binding:"required"`
	IsFavorite *bool `form:"isFavorite"`
}
