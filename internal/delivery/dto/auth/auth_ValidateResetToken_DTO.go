package dto

type Auth_ValidateResetToken_Request struct {
	Token string `form:"token" binding:"required"`
}
