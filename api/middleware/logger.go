package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Log() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// Additional fields
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method
		userAgent := c.Request.UserAgent()
		responseSize := c.Writer.Size()

		// Log the standard fields (like the default logger)
		fmt.Printf("[GIN] %v | %3d | %v | %s | %-7s %s | %d bytes | User-Agent: %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			responseSize,
			userAgent,
		)
	}
}
