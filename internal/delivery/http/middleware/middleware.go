package middleware

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func GeneralMiddleware() gin.HandlerFunc {
	reactFE := os.Getenv("URL_REACT")
	return cors.New(cors.Config{
		AllowOrigins:     []string{reactFE},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
	})
}
