package api

import (
	"io"
	"net/http"
	"strings"
	"github.com/gin-gonic/gin"
)

// Функция обработки запросов на сокращение URL
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	// Чтение тела запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Если произошла ошибка при чтении тела запроса, отправляем статус ошибки сервера и сообщение об ошибке
		c.String(http.StatusInternalServerError, "Ошибка при чтении тела запроса", http.StatusInternalServerError)
		return
	}
	// Удаление лишних пробелов в начале и конце строки
	URLtoBody := strings.TrimSpace(string(body))

	// Получение сокращенного URL с помощью сервиса структуры данных
	shortURL := s.StructService.GetShortURL(URLtoBody)

	// Установка типа содержимого ответа и отправка сокращенного URL в ответе
	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, shortURL)
}

// Функция обработки запросов на переадресацию на оригинальный URL
func (s *RestAPI) RedirectToOriginalURLHandler(c *gin.Context) {
	// Получение идентификатора сокращенного URL из параметра запроса
	shortID := c.Param("id")
	// Получение оригинального URL с помощью сервиса структуры данных
	originalURL, exists := s.StructService.GetOriginalURL(shortID)
	if !exists {
		// Если оригинального URL не найдено, отправляем статус временной переадресации и сообщение об ошибке
		c.String(http.StatusTemporaryRedirect, "URL не найден")
		return
	}
	// Установка заголовка "Location" для переадресации на оригинальный URL и отправка статуса временной переадресации и оригинального URL в ответе
	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, originalURL)
}
