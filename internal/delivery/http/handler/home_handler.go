package handler

import (
	dto "be-knowledge/internal/delivery/dto/home"
	"be-knowledge/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	homeService usecases.HomeService
}

func NewHomeHandler(homeService usecases.HomeService) *HomeHandler {
	return &HomeHandler{homeService}
}

func (h *HomeHandler) GetHistoryChat(c *gin.Context) {
	// namaEndpoint := "GetTopicById"
	// Tracelog.TopicLog("Mulai proses GetTopicById", namaEndpoint)

	var req dto.Home_GetHistoryChat_Request

	// Bind query param ke struct
	if err := c.ShouldBindQuery(&req); err != nil {
		//Tracelog.TopicLog("Request tidak valid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Panggil service dengan ID yang sudah tervalidasi
	topic, err := h.homeService.GetHistoryChat(req)
	if err != nil {
		//Tracelog.TopicLog("GetTopicById gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	//Tracelog.TopicLog("GetTopicById berhasil", namaEndpoint)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    topic,
	})
}
