package middleware

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"time"
)

// LoggerMiddleware создает middleware для логирования запросов с использованием zap.Logger.
func LoggerMiddleware(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now() // Запоминаем время начала обработки запроса
		c.Next()            // Передаем управление следующему обработчику

		duration := time.Since(start) // Вычисляем продолжительность обработки
		statusCode := c.Writer.Status() // Получаем код статуса ответа
		contentLength := int64(c.Writer.Size()) // Получаем размер содержимого ответа

		// Логируем информацию о запросе
		logger.Infow("Request",
			"method", c.Request.Method,        // Метод HTTP запроса
			"path", c.Request.URL.Path,        // Путь запроса
			"duration", duration,              // Продолжительность обработки
			"statusCode", statusCode,          // Код статуса ответа
			"contentLength", contentLength,    // Длина содержимого ответа
		)
	}
}
