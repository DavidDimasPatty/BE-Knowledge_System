package usecases

import (
	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
	"errors"
	"fmt"
)

type TopicService interface {
	GetTopicById(id int) (*entities.Topic, error)
	GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error)
	GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error)
}

type topicService struct {
	repo repository.TopicRepository
}

func NewTopicService(repo repository.TopicRepository) TopicService {
	return &topicService{repo}
}

func (s *topicService) GetTopicById(id int) (*entities.Topic, error) {
	topic, err := s.repo.GetTopicById(id)
	if err != nil {
		return nil, errors.New("topic tidak ditemukan")
	}
	return topic, nil
}

func (s *topicService) GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error) {
	topics, err := s.repo.GetAllTopicUser(username, isFavorite, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data topic: %w", err)
	}
	return topics, nil
}

func (s *topicService) GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error) {
	topics, err := s.repo.GetAllTopicUserByidCategories(username, idCategories)
	if err != nil {
		return nil, errors.New("gagal mengambil data topic")
	}
	return topics, nil
}
