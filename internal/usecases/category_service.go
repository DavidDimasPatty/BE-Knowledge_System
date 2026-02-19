package usecases

import (
	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
	Tracelog "be-knowledge/internal/tracelog"
	"fmt"
)

type CategoryService interface {
	GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error)
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func val[T any](p *T) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func (s *categoryService) GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error) {
	namaEndpoint := "GetAllCategoryUser"

	Tracelog.CategoryLog("Mulai proses usecase GetAllCategoryUser", namaEndpoint)

	Tracelog.CategoryLog(
		fmt.Sprintf("Parameter -> username: %s, search: %v, page: %v, limit: %v",
			username, val(search), val(page), val(limit)),
		namaEndpoint,
	)

	categories, err := s.repo.GetAllCategoryUser(username, search, page, limit)
	if err != nil {
		Tracelog.CategoryLog("Gagal mengambil data dari repository: "+err.Error(), namaEndpoint)
		return nil, fmt.Errorf("gagal mengambil data category: %w", err)
	}

	Tracelog.CategoryLog(
		fmt.Sprintf("Berhasil mengambil %d data category", len(categories)),
		namaEndpoint,
	)

	Tracelog.CategoryLog("Selesai proses usecase GetAllCategoryUser", namaEndpoint)

	return categories, nil
}
