package router

import (
	"be-knowledge/internal/delivery/http/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *handler.AuthHandler, userManagementHandler *handler.UserManagementHandler, dokumenManagementHandler *handler.DokumenManagementHandler, topicHandler *handler.TopicHandler, websocketHandler *handler.WebSocketHandler, categoryHandler *handler.CategoryHandler) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	registerPingRoutes(r)
	registerAuthRoutes(r, authHandler)
	registerGetAllUserRoutes(r, userManagementHandler)
	registerEditUserGet(r, userManagementHandler)
	registerAddUserRoutes(r, userManagementHandler)
	registerEditUserRoutes(r, userManagementHandler)
	registerDeleteUserRoutes(r, userManagementHandler)
	registerchangeStatusUserRoutes(r, userManagementHandler)
	registerGetAllDokumenRoutes(r, dokumenManagementHandler)
	registerEditDokumenGet(r, dokumenManagementHandler)
	registerAddDokumenRoutes(r, dokumenManagementHandler)
	registerEditDokumenRoutes(r, dokumenManagementHandler)
	registerDeleteDokumenRoutes(r, dokumenManagementHandler)
	registerDownloadDokumenRoutes(r, dokumenManagementHandler)
	registerTopicRoutes(r, topicHandler)
	registerWSChat(r, websocketHandler)
	registerCategoryRoutes(r, categoryHandler)
	return r
}

// Auth Router
func registerPingRoutes(r *gin.Engine) {
	r.GET("/ping", handler.Ping)
}

func registerAuthRoutes(r *gin.Engine, authHandler *handler.AuthHandler) {
	r.POST("/login", authHandler.Login)
}

// UserManagement Router
func registerGetAllUserRoutes(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.GET("/getAllUser", userManagementHandler.GetAllUser)
}
func registerEditUserGet(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.GET("/editUserGet", userManagementHandler.EditUserGet)
}
func registerAddUserRoutes(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.POST("/addUser", userManagementHandler.AddUser)
}
func registerEditUserRoutes(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.POST("/editUser", userManagementHandler.EditUser)
}
func registerDeleteUserRoutes(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.POST("/deleteUser", userManagementHandler.DeleteUser)
}
func registerchangeStatusUserRoutes(r *gin.Engine, userManagementHandler *handler.UserManagementHandler) {
	r.POST("/changeStatusUser", userManagementHandler.ChangeStatusUser)
}

// Dokumen Router
func registerGetAllDokumenRoutes(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.GET("/getAllDokumen", dokumenManagementHandler.GetAllDokumen)
}
func registerEditDokumenGet(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.GET("/editDokumenGet", dokumenManagementHandler.EditDokumenGet)
}
func registerAddDokumenRoutes(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.POST("/addDokumen", dokumenManagementHandler.AddDokumen)
}
func registerEditDokumenRoutes(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.POST("/editDokumen", dokumenManagementHandler.EditDokumen)
}
func registerDeleteDokumenRoutes(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.POST("/deleteDokumen", dokumenManagementHandler.DeleteDokumen)
}
func registerDownloadDokumenRoutes(r *gin.Engine, dokumenManagementHandler *handler.DokumenManagementHandler) {
	r.POST("/downloadDokumen", dokumenManagementHandler.DownloadDokumen)
}

// Topic Router
func registerTopicRoutes(r *gin.Engine, topicHandler *handler.TopicHandler) {
	r.GET("/getTopicById", topicHandler.GetTopicById)
	r.GET("/getAllTopicUser", topicHandler.GetAllTopicUser)
	r.GET("/getAllTopicUserByidCategories", topicHandler.GetAllTopicUserByidCategories)
}

// Web Socket
func registerWSChat(r *gin.Engine, WebSocketHandler *handler.WebSocketHandler) {
	r.GET("/ws", WebSocketHandler.Handle)
}

// Category Router
func registerCategoryRoutes(r *gin.Engine, categoryHandler *handler.CategoryHandler) {
	r.GET("/getAllCategoryUser", categoryHandler.GetAllCategoryUser)
}