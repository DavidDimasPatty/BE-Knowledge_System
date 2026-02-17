package handler

import (
	dto "be-knowledge/internal/delivery/dto/category"
	Tracelog "be-knowledge/internal/tracelog"
	"be-knowledge/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	categoryService usecases.CategoryService
}

func NewCategoryHandler(categoryService usecases.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService}
}

func (h *CategoryHandler) GetAllCategoryUser(c *gin.Context) {
	namaEndpoint := "GetAllCategoryUser"
	Tracelog.CategoryLog("Mulai proses GetAllCategoryUser", namaEndpoint)

	var req dto.Category_GetAllCategoryUser_Request

	if err := c.ShouldBindQuery(&req); err != nil {
		Tracelog.CategoryLog("Request tidak valid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	categories, err := h.categoryService.GetAllCategoryUser(req.Username, req.Search, req.Page, req.Limit)
	if err != nil {
		Tracelog.CategoryLog("GetAllCategoryUser gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Tracelog.CategoryLog("GetAllCategoryUser berhasil", namaEndpoint)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    categories,
	})
}
