package entities

import "time"

type User struct {
	ID              int        `db:"id"`
	Username        string     `db:"username"`
	Password        string     `db:"password"`
	Email           string     `db:"email"`
	NoTelp          string     `db:"noTelp"`
	Nama            string     `db:"nama"`
	Roles           int        `db:"roles"`
	Status          string     `db:"status"`
	LastLogin       *time.Time `db:"lastLogin"`
	OldPassword     string     `db:"oldPassword"`
	PasswordExpired *time.Time `db:"passwordExpired"`
	AddTime         *time.Time `db:"addTime"`
	UpdTime         *time.Time `db:"updTime"`
	AddId           string     `db:"addId"`
	UpdId           string     `db:"updId"`
}
