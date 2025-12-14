package usecases

import (
	"be-knowledge/internal/entities"
	"be-knowledge/internal/repository"
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

func (s *categoryService) GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error) {
	categories, err := s.repo.GetAllCategoryUser(username, search, page, limit)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil data category: %w", err)
	}
	return categories, nil
}
