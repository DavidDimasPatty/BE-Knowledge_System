package handler

import (
	dto "be-knowledge/internal/delivery/dto/auth"
	Tracelog "be-knowledge/internal/tracelog"
	"be-knowledge/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService usecases.UserService
}

func NewAuthHandler(userService usecases.UserService) *AuthHandler {
	return &AuthHandler{userService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var namaEndpoint = "Login"

	Tracelog.AuthLog("Mulai proses Login", namaEndpoint)

	var req dto.Auth_Login_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		Tracelog.AuthLog("Request tidak valid", namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		Tracelog.AuthLog("Login gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Tracelog.AuthLog("Login berhasil", namaEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"roleId":   user.Roles,
			"roleName": user.RoleName,
			"nama":     user.Nama,
			"email":    user.Email,
			"noTelp":   user.NoTelp,
		},
	})
}

func (h *AuthHandler) EditPassword(c *gin.Context) {
	var namaEndpoint = "EditPassword"

	Tracelog.AuthLog("Mulai proses EditPassword", namaEndpoint)

	var req dto.Auth_EditPassword_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		Tracelog.AuthLog("Request tidak valid", namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userService.EditPassword(req.Username, req.NewPassword, req.OldPassword)
	if err != nil {
		Tracelog.AuthLog("EditPassword gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Tracelog.AuthLog("EditPassword berhasil", namaEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (h *AuthHandler) SendEmailResetPassword(c *gin.Context) {
	var namaEndpoint = "SendEmailResetPassword"

	Tracelog.AuthLog("Mulai proses SendEmailResetPassword", namaEndpoint)

	var req dto.Auth_SendEmailResetPassword_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		Tracelog.AuthLog("Request tidak valid", namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userService.SendEmailResetPassword(req.Email)
	if err != nil {
		Tracelog.AuthLog("Gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Tracelog.AuthLog("Berhasil", namaEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (h *AuthHandler) ValidateResetToken(c *gin.Context) {
	var namaEndpoint = "ValidateResetToken"

	Tracelog.AuthLog("Mulai proses ValidateResetToken", namaEndpoint)

	var req dto.Auth_ValidateResetToken_Request

	if err := c.ShouldBindQuery(&req); err != nil {
		Tracelog.AuthLog("Token tidak valid", namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "token tidak valid"})
		return
	}

	err := h.userService.ValidateResetToken(req.Token)
	if err != nil {
		Tracelog.AuthLog("Token invalid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Tracelog.AuthLog("Token valid", namaEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var namaEndpoint = "ResetPassword"

	Tracelog.AuthLog("Mulai proses ResetPassword", namaEndpoint)

	var req dto.Auth_ResetPassword_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		Tracelog.AuthLog("Request tidak valid", namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userService.ResetPassword(req.Token, req.NewPassword)
	if err != nil {
		Tracelog.AuthLog("ResetPassword gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	Tracelog.AuthLog("ResetPassword berhasil", namaEndpoint)

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
	})
}