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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func (h *TopicHandler) GetAllTopicUser(c *gin.Context) {
	namaEndpoint := "GetAllTopicUser"
	Tracelog.TopicLog("Mulai proses GetAllTopicUser", namaEndpoint)

	var req dto.Topic_GetAllTopicUser_Request

	// Bind query param ke struct
	if err := c.ShouldBindQuery(&req); err != nil {
		Tracelog.TopicLog("Request tidak valid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topics, err := h.topicService.GetAllTopicUser(req.Username)
	if err != nil {
		Tracelog.TopicLog("GetAllTopicUser gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Tracelog.TopicLog("GetAllTopicUser berhasil", namaEndpoint)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    topics,
	})
}

func (h *TopicHandler) GetAllTopicUserByidCategories(c *gin.Context) {
	namaEndpoint := "GetAllTopicUserByidCategories"
	Tracelog.TopicLog("Mulai proses GetAllTopicUserByidCategories", namaEndpoint)

	var req dto.Topic_GetAllTopicUserByidCategories_Request

	// Bind query param ke struct
	if err := c.ShouldBindQuery(&req); err != nil {
		Tracelog.TopicLog("Request tidak valid: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	topics, err := h.topicService.GetAllTopicUserByidCategories(req.Username, req.IdCategories)
	if err != nil {
		Tracelog.TopicLog("GetAllTopicUserByidCategories gagal: "+err.Error(), namaEndpoint)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	Tracelog.TopicLog("GetAllTopicUserByidCategories berhasil", namaEndpoint)
	c.JSON(http.StatusOK, gin.H{
		"message": "Success",
		"data":    topics,
	})
}