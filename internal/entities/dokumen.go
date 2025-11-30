package entities

import "time"

type Dokumen struct {
	ID      int        `db:"id"`
	Link    string     `db:"link"`
	Judul   string     `db:"judul"`
	AddTime *time.Time `db:"ADDTIME"`
	UpdTime *time.Time `db:"updTime"`
	AddId   *string    `db:"addId"`
	UpdId   *string    `db:"updId"`
}
