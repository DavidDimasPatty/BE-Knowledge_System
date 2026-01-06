package dto

type Auth_EditPassword_Request struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"NewPassword" binding:"required"`
	OldPassword string `json:"OldPassword" binding:"required"`
}
