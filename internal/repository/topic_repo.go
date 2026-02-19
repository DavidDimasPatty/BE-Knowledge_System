package repository

import (
	"be-knowledge/internal/entities"
	Tracelog "be-knowledge/internal/tracelog"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type TopicRepository interface {
	GetTopicById(id int) (*entities.Topic, error)
	GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error)
	GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error)
	EditFavoriteTopic(username string, idTopic int, like int) error
}

type topicRepository struct {
	db *sqlx.DB
}

func NewTopicRepository(db *sqlx.DB) TopicRepository {
	return &topicRepository{db}
}

func (r *topicRepository) GetTopicById(id int) (*entities.Topic, error) {
	namaEndpoint := "GetTopicById"
	Tracelog.TopicLog("Mulai proses repository GetTopicById", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> id: %d", id),
		namaEndpoint,
	)

	topic := entities.Topic{}

	query := `
		SELECT 
			*
		FROM topic
		WHERE id = ? LIMIT 1
	`

	Tracelog.TopicLog(
		fmt.Sprintf("Query: %s | Params: [%d]", query, id),
		namaEndpoint,
	)

	err := r.db.Get(&topic, query, id)
	if err != nil {
		Tracelog.TopicLog("Query gagal: "+err.Error(), namaEndpoint)
		return nil, err
	}

	Tracelog.TopicLog("Berhasil mengambil topic", namaEndpoint)
	Tracelog.TopicLog("Selesai proses repository GetTopicById", namaEndpoint)

	return &topic, nil
}

func (r *topicRepository) GetAllTopicUser(username string, isFavorite *bool, search *string, page *int, limit *int) ([]entities.Topic, error) {
	namaEndpoint := "GetAllTopicUser"
	Tracelog.TopicLog("Mulai proses repository GetAllTopicUser", namaEndpoint)

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

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, isFavorite: %v, search: %v, page: %d, limit: %d, offset: %d",
			username,
			val(isFavorite),
			val(search),
			pageVal, limitVal, offset,
		),
		namaEndpoint,
	)

	baseQuery := `
        SELECT 
            t.*,b.category 
        FROM topic t 
        LEFT JOIN usertopicfavorite uf 
            ON uf.idTopic = t.id 
		JOIN categories b on t.idCategories=b.id 
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

	Tracelog.TopicLog(
		fmt.Sprintf("Query: %s | Params: %v", baseQuery, params),
		namaEndpoint,
	)

	// Execute
	err := r.db.Select(&topics, baseQuery, params...)
	if err != nil {
		Tracelog.TopicLog("Query gagal: "+err.Error(), namaEndpoint)
		return nil, err
	}

	Tracelog.TopicLog(
		fmt.Sprintf("Berhasil mengambil %d data topic", len(topics)),
		namaEndpoint,
	)

	Tracelog.TopicLog("Selesai proses repository GetAllTopicUser", namaEndpoint)

	return topics, nil
}

func (r *topicRepository) GetAllTopicUserByidCategories(username string, idCategories int) ([]entities.Topic, error) {
	namaEndpoint := "GetAllTopicUserByidCategories"
	Tracelog.TopicLog("Mulai proses repository GetAllTopicUserByidCategories", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, idCategories: %d",
			username, idCategories),
		namaEndpoint,
	)

	topics := []entities.Topic{}

	query := `
		SELECT 
			*
		FROM topic
		WHERE addId = ? AND idCategories = ?
		ORDER BY ADDTIME DESC
	`

	Tracelog.TopicLog(
		fmt.Sprintf("Query: %s | Params: [%s, %d]",
			query, username, idCategories),
		namaEndpoint,
	)

	err := r.db.Select(&topics, query, username, idCategories)
	if err != nil {
		Tracelog.TopicLog("Query gagal: "+err.Error(), namaEndpoint)
		return nil, err
	}

	Tracelog.TopicLog(
		fmt.Sprintf("Berhasil mengambil %d data topic", len(topics)),
		namaEndpoint,
	)

	Tracelog.TopicLog("Selesai proses repository GetAllTopicUserByidCategories", namaEndpoint)

	return topics, nil
}

func (r *topicRepository) EditFavoriteTopic(username string, idTopic int, like int) error {
	namaEndpoint := "EditFavoriteTopic"
	Tracelog.TopicLog("Mulai proses repository EditFavoriteTopic", namaEndpoint)

	Tracelog.TopicLog(
		fmt.Sprintf("Parameter -> username: %s, idTopic: %d, like: %d",
			username, idTopic, like),
		namaEndpoint,
	)

	var err error
	if like == 0 {
		query := "DELETE FROM usertopicfavorite WHERE idTopic = ? and addId=?"
		Tracelog.TopicLog(
			fmt.Sprintf("Query: %s | Params: [%d, %s]",
				query, idTopic, username),
			namaEndpoint,
		)
		_, err = r.db.Exec(query, idTopic, username)
	} else {
		query := "INSERT INTO usertopicfavorite values(?,now(),?)"
		Tracelog.TopicLog(
			fmt.Sprintf("Query: %s | Params: [%d, %s]",
				query, idTopic, username),
			namaEndpoint,
		)
		_, err = r.db.Exec(query, idTopic, username)
	}

	if err != nil {
		Tracelog.TopicLog("Query gagal: "+err.Error(), namaEndpoint)
		return err
	}

	Tracelog.TopicLog("Berhasil update favorite topic", namaEndpoint)
	Tracelog.TopicLog("Selesai proses repository EditFavoriteTopic", namaEndpoint)

	return err
}
