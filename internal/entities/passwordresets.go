package entities

import "time"

type PasswordResets struct {
	ID          int        `db:"id"`
	User_Id     int        `db:"user_id"`
	Token       *string    `db:"token"`
	ExpiredDate *time.Time `db:"expiredDate"`
	AddTime     *time.Time `db:"addTime"`
	IsReset     *string    `db:"isReset"`
}
