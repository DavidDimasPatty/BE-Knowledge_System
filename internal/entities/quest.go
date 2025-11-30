package entities

import "time"

type Quest struct {
	ID      *int       `db:"id"`
	Isi     *string    `db:"isi"`
	IdChat  *int       `db:"int"`
	AddTime *time.Time `db:"ADDTIME"`
}
