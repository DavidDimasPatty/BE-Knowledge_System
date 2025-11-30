package entities

import "time"

type Answer struct {
	ID      int       `db:"id"`
	Isi     string    `db:"isi"`
	Step    int       `db:"step"`
	IsLike  int       `db:"isLike"`
	AddTime time.Time `db:"ADDTIME"`
	IdQuest int       `db:"idQuest"`
	Refer   string    `db:"refer"`
}
