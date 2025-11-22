package router

import (
	"github.com/gin-gonic/gin"
	"example/be-knowledge/internal/delivery/http/handler"
)

func NewRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", handler.Ping)

	return r
}