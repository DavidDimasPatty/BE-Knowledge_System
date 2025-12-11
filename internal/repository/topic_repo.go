package repository

import (
	"be-knowledge/internal/entities"

	"github.com/jmoiron/sqlx"
)

type TopicRepository interface {
	GetTopicById(id int) (*entities.Topic, error)
	GetAllTopicUser(username string, isFavorite *bool) ([]entities.Topic, error)
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

func (r *topicRepository) GetAllTopicUser(username string, isFavorite *bool) ([]entities.Topic, error) {
    topics := []entities.Topic{}

    baseQuery := `
        SELECT 
            t.* 
        FROM topic t 
        LEFT JOIN usertopicfavorite uf 
            ON uf.idTopic = t.id 
        WHERE t.addId = ? 
    `

    // Filter berdasarkan IsFavorite
    if isFavorite != nil {
        if *isFavorite {
            // hanya favorite
            baseQuery += " AND uf.idTopic IS NOT NULL"
        } else {
            // hanya non-favorite
            baseQuery += " AND uf.idTopic IS NULL"
        }
    }

    // Urutkan terbaru
    baseQuery += " ORDER BY t.addTime DESC"

    err := r.db.Select(&topics, baseQuery, username)
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
