package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				log.Printf("Request error: %v", e.Err)
			}
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Terjadi kesalahan server",
				"data":    nil,
			})
		}
	}
}