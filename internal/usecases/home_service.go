package usecases

import (
	dto "be-knowledge/internal/delivery/dto/home"
	"be-knowledge/internal/repository"
	Tracelog "be-knowledge/internal/tracelog"
	"errors"
	"fmt"
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
	namaEndpoint := "GetHistoryChat"

	Tracelog.HomeLog("Mulai proses usecase GetHistoryChat", namaEndpoint)

	Tracelog.HomeLog(
		fmt.Sprintf("Parameter -> %+v", data),
		namaEndpoint,
	)

	res, err := s.repo.GetHistoryChat(data)
	if err != nil {
		Tracelog.HomeLog("Gagal mengambil history chat: "+err.Error(), namaEndpoint)
		return nil, errors.New("topic tidak ditemukan")
	}

	Tracelog.HomeLog("Berhasil mengambil history chat", namaEndpoint)
	Tracelog.HomeLog("Selesai proses usecase GetHistoryChat", namaEndpoint)

	return res, nil
}
