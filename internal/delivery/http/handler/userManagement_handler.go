package handler

import (
	dto "be-knowledge/internal/delivery/dto/userManagement"
	"be-knowledge/internal/usecases"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserManagementHandler struct {
	userManagementService usecases.UserManagementService
}

func NewUserManagementHandler(userManagementService usecases.UserManagementService) *UserManagementHandler {
	return &UserManagementHandler{userManagementService}
}

func (h *UserManagementHandler) GetAllUser(c *gin.Context) {

	user, err := h.userManagementService.GetAllUser()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"user": gin.H{
			"data": user.Data,
		},
	})
}

func (h *UserManagementHandler) EditUserGet(c *gin.Context) {
	idStr := c.Query("id")
	if idStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a number"})
		return
	}

	user, err := h.userManagementService.EditUserGet(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    user,
	})
}

func (h *UserManagementHandler) AddUser(c *gin.Context) {
	var req *dto.UserManagement_AddUser_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userManagementService.AddUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    "null",
	})
}

func (h *UserManagementHandler) EditUser(c *gin.Context) {
	var req *dto.UserManagement_EditUser_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userManagementService.EditUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    "null",
	})
}

func (h *UserManagementHandler) DeleteUser(c *gin.Context) {
	var req *dto.UserManagement_DeleteUser_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.userManagementService.DeleteUser(req.Id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    "null",
	})
}
