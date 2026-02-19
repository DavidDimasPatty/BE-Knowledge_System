package repository

import (
	"be-knowledge/internal/entities"
	Tracelog "be-knowledge/internal/tracelog"
	"fmt"

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

func val[T any](p *T) interface{} {
	if p == nil {
		return nil
	}
	return *p
}

func (r *categoryRepository) GetAllCategoryUser(username string, search *string, page *int, limit *int) ([]entities.Category, error) {
	namaEndpoint := "GetAllCategoryUser"
	Tracelog.CategoryLog("Mulai proses repository GetAllCategoryUser", namaEndpoint)

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

	Tracelog.CategoryLog(
		fmt.Sprintf("Parameter -> username: %s, search: %v, page: %d, limit: %d, offset: %d",
			username, val(search), pageVal, limitVal, offset),
		namaEndpoint,
	)

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
	baseQuery += " ORDER BY c.addTime DESC LIMIT ? OFFSET ?"
	params = append(params, limitVal, offset)

	Tracelog.CategoryLog(
		fmt.Sprintf("Final Query: %s | Params: %v", baseQuery, params),
		namaEndpoint,
	)

	// Execute
	err := r.db.Select(&categories, baseQuery, params...)
	if err != nil {
		Tracelog.CategoryLog("Query gagal: "+err.Error(), namaEndpoint)
		return nil, err
	}

	Tracelog.CategoryLog(
		fmt.Sprintf("Berhasil mengambil %d data category", len(categories)),
		namaEndpoint,
	)

	Tracelog.CategoryLog("Selesai proses repository GetAllCategoryUser", namaEndpoint)

	return categories, nil
}
