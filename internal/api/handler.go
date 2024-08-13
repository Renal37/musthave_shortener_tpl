package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// Структура для запроса URL
type Request struct {
	URL string `json:"url"`
}

// Структура для ответа с сокращённым URL
type Response struct {
	Result string `json:"result"`
}

// Структура для запросов с идентификатором корреляции и оригинальным URL
type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Структура для ответов с идентификатором корреляции и сокращённым URL
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Обработчик для сокращения URL
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated
	body, err := io.ReadAll(c.Request.Body) // Чтение тела запроса
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}
	url := strings.TrimSpace(string(body)) // Удаление лишних пробелов
	shortURL, err := s.StructService.Set(url) // Сокращение URL
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err) // Проверка, существует ли уже сокращённый URL
		if err != nil {
			c.String(http.StatusInternalServerError, "URL не может быть сокращён", http.StatusInternalServerError)
			return
		}
		httpStatus = http.StatusConflict
	}
	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

// Обработчик для сокращения URL с использованием JSON
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body) // Декодирование JSON из тела запроса
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	url := strings.TrimSpace(decoderBody.URL) // Удаление лишних пробелов
	shortURL, err := s.StructService.Set(url) // Сокращение URL
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err) // Проверка, существует ли уже сокращённый URL
		if err != nil {
			errorMessage := map[string]interface{}{
				"message": "URL не может быть сокращён",
				"code":    http.StatusInternalServerError,
			}
			answer, _ := json.Marshal(errorMessage)
			c.Data(http.StatusInternalServerError, "application/json", answer)
			return
		}
		httpStatus = http.StatusConflict
	}

	StructPerformance := Response{Result: shortURL}
	respJSON, err := json.Marshal(StructPerformance)
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Обработчик для перенаправления на оригинальный URL
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	shortID := c.Param("id")
	originalURL, exists := s.StructService.Get(shortID) // Получение оригинального URL по сокращённому
	if !exists {
		c.String(http.StatusTemporaryRedirect, "URL не найден")
		return
	}
	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, originalURL)
}

// Обработчик для сокращения нескольких URL с использованием JSON
func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body) // Декодирование JSON из тела запроса
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL) // Удаление лишних пробелов
		shortURL, err := s.StructService.Set(url) // Сокращение URL
		if err != nil {
			shortURL, err = s.StructService.GetExistURL(url, err) // Проверка, существует ли уже сокращённый URL
			if err != nil {
				errorMessage := map[string]interface{}{
					"message": "URL не может быть сокращён",
					"code":    http.StatusInternalServerError,
				}
				answer, _ := json.Marshal(errorMessage)
				c.Data(http.StatusInternalServerError, "application/json", answer)
				return
			}
			httpStatus = http.StatusConflict
		}
		urlResponse := ResponseBodyURLs{
			req.CorrelationID,
			shortURL,
		}
		URLResponses = append(URLResponses, urlResponse)
	}
	respJSON, err := json.Marshal(URLResponses)
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMessage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Обработчик для проверки доступности сервиса
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.StructService.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}
