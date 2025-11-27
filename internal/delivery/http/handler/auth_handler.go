package handler

import (
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

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
		"message": "login berhasil",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"roles":    user.Roles,
		},
	})
}
