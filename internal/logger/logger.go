package logger

import (
	"time"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var Log = zap.NewNop()

func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl
	return nil
}

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)
		statusCode := c.Writer.Status()
		contentLength := int64(c.Writer.Size())

		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path1", c.Request.URL.Path),
			zap.Duration("duration2", duration),
		)

		logger.Info("Response",
			zap.Int("statusCode", statusCode),
			zap.Int64("contentLength", contentLength),
		)
	}
}
