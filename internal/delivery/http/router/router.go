package router

import (
	"be-knowledge/internal/delivery/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler) *gin.Engine {
	r := gin.Default()

	registerPingRoutes(r)
	registerAuthRoutes(r, authHandler)

	return r
}

func registerPingRoutes(r *gin.Engine) {
	r.GET("/ping", handler.Ping)
}

func registerAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	r.POST("/login", authHandler.Login)
}

// func getAllUserRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
// 	r.POST("/login", authHandler.Login)
// }
