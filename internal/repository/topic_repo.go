package repository

import (
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"
)

type TopicRepository interface {
	GetTopicById(id int) (*entities.Topic, error)
	GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error)
	GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error)
}

type topicRepository struct {
	db *sqlx.DB
}

func NewTopicRepository(db *sqlx.DB) TopicRepository {
	return &topicRepository{db}
}

func (r *topicRepository) GetTopicById(id int) (*entities.Topic, error) {
	topic := entities.Topic{}

	query := `
		SELECT 
			*
		FROM topic
		WHERE id = ? LIMIT 1
	`

	err := r.db.Get(&topic, query, id)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func (r *topicRepository) GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error) {
    topics := []entities.Topic{}

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
        SELECT 
            t.* 
        FROM topic t 
        LEFT JOIN usertopicfavorite uf 
            ON uf.idTopic = t.id 
        WHERE t.addId = ?
    `
    params := []interface{}{username}

    // Filter berdasarkan IsFavorite
    if isFavorite != nil {
        if *isFavorite {
            baseQuery += " AND uf.idTopic IS NOT NULL"
        } else {
            baseQuery += " AND uf.idTopic IS NULL"
        }
    }

    // Filter search
    if search != nil && *search != "" {
        baseQuery += " AND (t.topic LIKE ? OR t.descriptions LIKE ?)"
        like := "%" + *search + "%"
        params = append(params, like, like)
    }

    // Order & Pagination
    baseQuery += " ORDER BY t.addTime DESC LIMIT ? OFFSET ?"
    params = append(params, limitVal, offset)

    // Execute
    err := r.db.Select(&topics, baseQuery, params...)
    if err != nil {
        return nil, err
    }

    return topics, nil
}

func (r *topicRepository) GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error) {
	topics := []entities.Topic{}

	query := `
		SELECT 
			*
		FROM topic
		WHERE addId = ? AND idCategories = ?
		ORDER BY ADDTIME DESC
	`

	err := r.db.Select(&topics, query, username, idCategories)
	if err != nil {
		return nil, err
	}

	return topics, nil
}
