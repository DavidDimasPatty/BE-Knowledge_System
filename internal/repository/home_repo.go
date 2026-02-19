package repository

import (
	dto "be-knowledge/internal/delivery/dto/home"
	Tracelog "be-knowledge/internal/tracelog"
	"fmt"

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

	namaEndpoint := "GetHistoryChat"
	Tracelog.HomeLog("Mulai proses repository GetHistoryChat", namaEndpoint)

	Tracelog.HomeLog(
		fmt.Sprintf("Parameter -> Username: %s, Category: %v, Topic: %v",
			data.Username, data.Category, data.Topic),
		namaEndpoint,
	)

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

	Tracelog.HomeLog(
		fmt.Sprintf("Query: %s | Params: [%s, %v, %v]",
			query,
			data.Username,
			data.Category,
			data.Topic,
		),
		namaEndpoint,
	)

	err := r.db.Select(&rows, query,
		data.Username,
		data.Category,
		data.Topic,
	)
	if err != nil {
		Tracelog.HomeLog("Query gagal: "+err.Error(), namaEndpoint)
		return nil, err
	}

	Tracelog.HomeLog(
		fmt.Sprintf("Berhasil mengambil %d baris history chat", len(rows)),
		namaEndpoint,
	)

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

	Tracelog.HomeLog("Selesai proses repository GetHistoryChat", namaEndpoint)

	return &resp, nil
}
