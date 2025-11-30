package handler

import (
	dto "be-knowledge/internal/delivery/dto/auth"
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
	var req dto.Auth_Login_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	user, err := h.userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"roleId":   user.Roles,
			"roleName": user.RoleName,
		},
	})
}
