package middleware

import (
	"time"

	"theb-backend/internal/logger"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware for logging HTTP requests
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		// Process request
		c.Next()

		// Log after request
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		fields := map[string]interface{}{
			"method":     method,
			"path":       path,
			"status":     statusCode,
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}

		if requestID := c.GetString("request_id"); requestID != "" {
			fields["request_id"] = requestID
		}

		if statusCode >= 500 {
			logger.Error("HTTP Request", fields)
		} else if statusCode >= 400 {
			logger.Warn("HTTP Request", fields)
		} else {
			logger.Info("HTTP Request", fields)
		}
	}
}
