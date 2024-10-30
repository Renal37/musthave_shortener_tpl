package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

func LoggerMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		statusCode := c.Writer.Status()
		contentLength := int64(c.Writer.Size())

		logger.Infow("Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"duration", duration,
			"statusCode", statusCode,
			"contentLength", contentLength,
		)
	}
}
