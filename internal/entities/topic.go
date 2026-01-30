package entities

import "time"

type Topic struct {
	ID           *int       `db:"id"`
	Topic        *string    `db:"topic"`
	Desctription *string    `db:"descriptions"`
	Category     *string    `db:"category"`
	IdCategories *int       `db:"idCategories"`
	AddTime      *time.Time `db:"ADDTIME"`
	UpdTime      *time.Time `db:"updTime"`
	AddId        *string    `db:"addId"`
	UpdId        *string    `db:"updId"`
}
