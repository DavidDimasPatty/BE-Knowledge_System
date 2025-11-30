package entities

type Chat struct {
	ID      int    `db:"id"`
	IdTopic string `db:"idTopic"`
}
