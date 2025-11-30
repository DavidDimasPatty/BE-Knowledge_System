package handler

import (
	dto "be-knowledge/internal/delivery/dto/dokumenManagement"
	"be-knowledge/internal/usecases"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type DokumenManagementHandler struct {
	dokumenManagementService usecases.DokumenManagementService
}

func NewDokumenManagementHandler(dokumenManagementService usecases.DokumenManagementService) *DokumenManagementHandler {
	return &DokumenManagementHandler{dokumenManagementService}
}

func (h *DokumenManagementHandler) GetAllDokumen(c *gin.Context) {

	dokumen, err := h.dokumenManagementService.GetAllDokumen()
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"dokumen": gin.H{
			"data": dokumen.Data,
		},
	})
}

func (h *DokumenManagementHandler) EditDokumenGet(c *gin.Context) {
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

	dokumen, err := h.dokumenManagementService.EditDokumenGet(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    dokumen,
	})
}

func (h *DokumenManagementHandler) AddDokumen(c *gin.Context) {
	var req *dto.DokumenManagement_AddDokumen_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.dokumenManagementService.AddDokumen(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    nil,
	})
}

func (h *DokumenManagementHandler) EditDokumen(c *gin.Context) {
	var req *dto.DokumenManagement_EditDokumen_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.dokumenManagementService.EditDokumen(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    nil,
	})
}

func (h *DokumenManagementHandler) DeleteDokumen(c *gin.Context) {
	var req *dto.DokumenManagement_DeleteDokumen_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	err := h.dokumenManagementService.DeleteDokumen(req.Id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    nil,
	})
}

func (h *DokumenManagementHandler) DownloadDokumen(c *gin.Context) {
	var req *dto.DokumenManagement_DeleteDokumen_Request

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	dok, fileBytes, err := h.dokumenManagementService.DownloadDokumen(req.Id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(dok.Link))
	c.Data(200, "application/octet-stream", fileBytes)
}
