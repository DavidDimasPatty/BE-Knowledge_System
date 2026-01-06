package dto

type Topic_EditFavoriteTopic_Request struct {
	IdTopic  int    `json:"idTopic"`
	Username string `json:"username"`
	Like     int    `json:"like"`
}
