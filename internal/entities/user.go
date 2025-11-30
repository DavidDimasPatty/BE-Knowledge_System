package entities

import "time"

type User struct {
	ID              int        `db:"id"`
	Username        string     `db:"username"`
	Password        string     `db:"PASSWORD"`
	Email           *string    `db:"email"`
	NoTelp          *string    `db:"noTelp"`
	Nama            *string    `db:"nama"`
	Roles           int        `db:"roles"`
	Status          *string    `db:"STATUS"`
	LastLogin       *time.Time `db:"lastLogin"`
	LoginCount      int       `db:"loginCount"`
	BlockDate       *time.Time `db:"blockDate"`
	OldPassword     *string    `db:"oldPassword"`
	PasswordExpired *time.Time `db:"passwordExpired"`
	AddTime         *time.Time `db:"ADDTIME"`
	UpdTime         *time.Time `db:"updTime"`
	AddId           *string    `db:"addId"`
	UpdId           *string    `db:"updId"`
	Divisi          int        `db:"divisi"`

	RoleName *string `db:"role_name"`
}
