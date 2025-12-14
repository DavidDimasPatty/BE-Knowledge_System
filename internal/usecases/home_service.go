package usecases

import (
	dto "be-knowledge/internal/delivery/dto/home"
	"be-knowledge/internal/repository"
	"errors"
)

type HomeService interface {
	GetHistoryChat(data dto.Home_GetHistoryChat_Request) (*dto.Home_GetHistoryChat_Response, error)
}

type homeService struct {
	repo repository.HomeRepository
}

func NewHomeService(repo repository.HomeRepository) HomeService {
	return &homeService{repo}
}

func (s *homeService) GetHistoryChat(data dto.Home_GetHistoryChat_Request) (*dto.Home_GetHistoryChat_Response, error) {
	res, err := s.repo.GetHistoryChat(data)
	if err != nil {
		return nil, errors.New("topic tidak ditemukan")
	}
	return res, nil
}
