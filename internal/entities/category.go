package entities

import "time"

type Category struct {
	ID       int       `db:"id"`
	Category string    `db:"category"`
	Icon     int       `db:"icon"`
	AddTime  time.Time `db:"ADDTIME"`
	UpdTime  time.Time `db:"updTime"`
	AddId    string    `db:"addId"`
	UpdId    string    `db:"updId"`
}
