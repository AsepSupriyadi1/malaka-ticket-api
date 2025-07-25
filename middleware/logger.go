package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func CustomLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		c.Next()
		duration := time.Since(start)
		status := c.Writer.Status()

		log.Printf("%s %s | %d | %v", method, path, status, duration)
	}
}
