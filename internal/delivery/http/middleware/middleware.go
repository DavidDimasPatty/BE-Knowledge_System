package middleware

import "github.com/gin-gonic/gin"

// Middleware kosong (placeholder)
// Tidak melakukan apa-apa, hanya melanjutkan request
func GeneralMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// bisa tambahkan log atau apa pun di sini
		c.Next()
	}
}
