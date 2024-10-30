package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// LoggerMiddleware возвращает middleware для логирования запросов и ответов с использованием zap.Logger
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Запоминаем время начала запроса
		start := time.Now()

		// Обработка запроса
		c.Next()

		// Вычисляем длительность запроса
		duration := time.Since(start)

		// Получаем статус-код ответа
		statusCode := c.Writer.Status()

		// Получаем длину содержимого ответа
		contentLength := c.Writer.Size()
		if contentLength < 0 {
			contentLength = 0
		}

		// Логируем информацию о запросе
		logger.Info("Request",
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), // Исправлено название поля на "path"
			zap.Duration("duration", duration),     // Исправлено название поля на "duration"
			zap.String("Response", ""),
		)

		// Логируем информацию об ответе
		logger.Info("Response",
			zap.Int("statusCode", statusCode),
			zap.Int("contentLength", contentLength), // Исправлено тип поля на zap.Int, так как Size() возвращает int
		)
	}
}