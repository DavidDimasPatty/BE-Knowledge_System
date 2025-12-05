package dto

import "be-knowledge/internal/entities"

type Topic_GetTopicById_Response struct {
	Data []entities.Topic `json:"data"`
}

type Topic_GetTopicById_Request struct {
	Id int `form:"id" binding:"required"`
}
