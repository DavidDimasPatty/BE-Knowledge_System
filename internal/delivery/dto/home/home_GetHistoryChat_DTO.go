package dto

type Home_GetHistoryChat_Response struct {
	User []ChatHistory `json:"user"`
	Bot  []ChatHistory `json:"bot"`
}

type Home_GetHistoryChat_Request struct {
	Username string `form:"username"`
	Topic    int    `form:"topic"`
	Category int    `form:"category"`
}

type ChatHistory struct {
	Isi  string
	Role string
}

type ChatHistoryRow struct {
	User string `db:"user"`
	Bot  string `db:"bot"`
}
