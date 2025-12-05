package handler

import (
	dto "be-knowledge/internal/delivery/dto/topic"
	Tracelog "be-knowledge/internal/tracelog"
	"be-knowledge/internal/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TopicHandler struct {
	topicService usecases.TopicService
}

func NewTopicHandler(topicService usecases.TopicService) *TopicHandler {
	return &TopicHandler{topicService}
}

func (h *TopicHandler) GetTopicById(c *gin.Context) {
	namaEndpoint := "GetTopicById"
	Tracelog.TopicLog("Mulai proses GetTopicById", namaEndpoint)

	var req dto.Topic_GetTopicById_Request

	// Bind query param ke struct
	if err := c.ShouldBindQuery(&req); err != nil {
		Tracelog.TopicLog("Request tidak valid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required & must be a number"})
		return
	}

	// Panggil service dengan ID yang sudah tervalidasi
	topic, err := h.topicService.GetTopicById(req.Id)
	if err != nil {
		Tracelog.TopicLog("GetTopicById gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	Tracelog.TopicLog("GetTopicById berhasil", namaEndpoint)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    topic,
	})
}

// func (h *TopicHandler) GetAllTopic(c *gin.Context) {
// 	topics, err := h.topicService.GetAllTopic()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Success",
// 		"data":    topics,
// 	})
// }
