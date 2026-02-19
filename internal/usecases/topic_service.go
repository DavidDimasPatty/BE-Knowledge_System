package usecases

import (
	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
	Tracelog "be-knowledge/internal/tracelog"
	"errors"
	"fmt"
)

type TopicService interface {
	GetTopicById(id int) (*entities.Topic, error)
	GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error)
	GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error)
	EditFavoriteTopic(username string, idTopic int, like int) error
}

type topicService struct {
	repo repository.TopicRepository
}

func NewTopicService(repo repository.TopicRepository) TopicService {
	return &topicService{repo}
}

func (s *topicService) GetTopicById(id int) (*entities.Topic, error) {
	namaEndpoint := "GetTopicById"
	Tracelog.TopicLog("Mulai proses usecase GetTopicById", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> id: %d", id),
		namaEndpoint,
	)

	topic, err := s.repo.GetTopicById(id)
	if err != nil {
		Tracelog.TopicLog("Topic tidak ditemukan: "+err.Error(), namaEndpoint)
		return nil, errors.New("topic tidak ditemukan")
	}

	Tracelog.TopicLog("Berhasil mengambil topic", namaEndpoint)
	Tracelog.TopicLog("Selesai proses usecase GetTopicById", namaEndpoint)
	return topic, nil
}

func (s *topicService) GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error) {
	namaEndpoint := "GetAllTopicUser"
	Tracelog.TopicLog("Mulai proses usecase GetAllTopicUser", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, isFavorite: %v, search: %v, page: %v, limit: %v",
			username,
			val(isFavorite),
			val(search),
			val(page),
			val(limit),
		),
		namaEndpoint,
	)

	topics, err := s.repo.GetAllTopicUser(username, isFavorite, search, page, limit)
	if err != nil {
		Tracelog.TopicLog("Gagal mengambil data topic: "+err.Error(), namaEndpoint)
		return nil, fmt.Errorf("gagal mengambil data topic: %w", err)
	}

	Tracelog.TopicLog(
		fmt.Sprintf("Berhasil mengambil %d data topic", len(topics)),
		namaEndpoint,
	)

	Tracelog.TopicLog("Selesai proses usecase GetAllTopicUser", namaEndpoint)

	return topics, nil
}

func (s *topicService) GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error) {
	namaEndpoint := "GetAllTopicUserByidCategories"
	Tracelog.TopicLog("Mulai proses usecase GetAllTopicUserByidCategories", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, idCategories: %d",
			username, idCategories),
		namaEndpoint,
	)

	topics, err := s.repo.GetAllTopicUserByidCategories(username, idCategories)
	if err != nil {
		Tracelog.TopicLog("Gagal mengambil data topic: "+err.Error(), namaEndpoint)
		return nil, errors.New("gagal mengambil data topic")
	}

	Tracelog.TopicLog(
		fmt.Sprintf("Berhasil mengambil %d data topic", len(topics)),
		namaEndpoint,
	)

	Tracelog.TopicLog("Selesai proses usecase GetAllTopicUserByidCategories", namaEndpoint)

	return topics, nil
}

func (s *topicService) EditFavoriteTopic(username string, idTopic int, like int) error {
	namaEndpoint := "EditFavoriteTopic"
	Tracelog.TopicLog("Mulai proses usecase EditFavoriteTopic", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, idTopic: %d, like: %d",
			username, idTopic, like),
		namaEndpoint,
	)

	err := s.repo.EditFavoriteTopic(username, idTopic, like)
	if err != nil {
		Tracelog.TopicLog("Gagal update favorite topic: "+err.Error(), namaEndpoint)
		return err
	}

	Tracelog.TopicLog("Berhasil update favorite topic", namaEndpoint)
	Tracelog.TopicLog("Selesai proses usecase EditFavoriteTopic", namaEndpoint)

	return nil
}
