package dto

type Auth_SendEmailResetPassword_Request struct {
	Email string `json:"email" binding:"required,email"`
}
