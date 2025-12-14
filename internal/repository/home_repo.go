package repository

import (
	dto "be-knowledge/internal/delivery/dto/home"

	"github.com/jmoiron/sqlx"
)

type HomeRepository interface {
	GetHistoryChat(data dto.Home_GetHistoryChat_Request) (*dto.Home_GetHistoryChat_Response, error)
}

type homeRepository struct {
	db *sqlx.DB
}

func NewHomeRepository(db *sqlx.DB) HomeRepository {
	return &homeRepository{db}
}

func (r *homeRepository) GetHistoryChat(
	data dto.Home_GetHistoryChat_Request,
) (*dto.Home_GetHistoryChat_Response, error) {

	rows := []dto.ChatHistoryRow{}

	query := `
	SELECT 
		a.isi AS user,
		e.isi AS bot
	FROM quest a
	JOIN chats b ON a.idChat = b.id
	JOIN topic c ON b.idTopic = c.id
	JOIN categories d ON c.idCategories = d.id
	JOIN answer e ON e.idQuest = a.id
	WHERE d.addId = ?
	  AND d.id = ?
	  AND c.id = ?
	`

	err := r.db.Select(&rows, query,
		data.Username,
		data.Category,
		data.Topic,
	)
	if err != nil {
		return nil, err
	}

	resp := dto.Home_GetHistoryChat_Response{
		User: []dto.ChatHistory{},
		Bot:  []dto.ChatHistory{},
	}

	for _, r := range rows {
		resp.User = append(resp.User, dto.ChatHistory{
			Isi:  r.User,
			Role: "user",
		})

		resp.Bot = append(resp.Bot, dto.ChatHistory{
			Isi:  r.Bot,
			Role: "bot",
		})
	}

	return &resp, nil
}
