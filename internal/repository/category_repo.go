package repository

import (
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"
)

type CategoryRepository interface {
	GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error)
}

type categoryRepository struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error) {
	categories := []entities.Category{}

	// Default values jika nil
	pageVal := 1
	limitVal := 20

	if page != nil && *page > 0 {
		pageVal = *page
	}
	if limit != nil && *limit > 0 {
		limitVal = *limit
	}

	offset := (pageVal - 1) * limitVal

	baseQuery := `
        SELECT DISTINCT
            c.*, i.nama AS nama_icon
        FROM topic t 
		LEFT JOIN categories c ON c.id = t.idCategories
		LEFT JOIN icon i ON i.id = c.icon
        WHERE t.addId = ?
    `
	params := []interface{}{username}

	// Filter search
	if search != nil && *search != "" {
		baseQuery += " AND (c.category LIKE ?)"
		like := "%" + *search + "%"
		params = append(params, like)
	}

	// Order & Pagination
	baseQuery += " ORDER BY t.addTime DESC LIMIT ? OFFSET ?"
	params = append(params, limitVal, offset)

	// Execute
	err := r.db.Select(&categories, baseQuery, params...)
	if err != nil {
		return nil, err
	}

	return categories, nil
}
