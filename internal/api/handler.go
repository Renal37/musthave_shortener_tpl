package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Структура для получения URL из тела запроса
type URLRequest struct {
	URL string `json:"url"`
}

// Структура для ответа с коротким URL
type URLResponse struct {
	Result string `json:"result"`
}

// Структура для тела запроса с множеством URL
type BulkURLRequest struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Структура для ответа с множеством коротких URL
type BulkURLResponse struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Обработчик сокращения URL из текстового тела запроса
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	var httpStatus = http.StatusCreated

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка при чтении тела запроса")
		return
	}

	userID, _ := c.Get("userID")
	url := strings.TrimSpace(string(body))

	shortURL, err := s.Shortener.Set(userID.(string), url)
	if err != nil {
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "Ошибка при сокращении URL")
			return
		}
		httpStatus = http.StatusConflict
	}

	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

// Обработчик сокращения URL в формате JSON
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var req URLRequest
	httpStatus := http.StatusCreated

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка при чтении тела запроса",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	userID, _ := c.Get("userID")
	url := strings.TrimSpace(req.URL)

	shortURL, err := s.Shortener.Set(userID.(string), url)
	if err != nil {
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Ошибка при сокращении URL",
				"code":    http.StatusInternalServerError,
			})
			return
		}
		httpStatus = http.StatusConflict
	}

	c.JSON(httpStatus, URLResponse{Result: shortURL})
}

// Обработчик сокращения множества URL (JSON формат)
func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var requests []BulkURLRequest
	httpStatus := http.StatusCreated

	if err := c.ShouldBindJSON(&requests); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка при чтении тела запроса",
			"code":    http.StatusInternalServerError,
		})
		return
	}

	userID, _ := c.Get("userID")

	var responses []BulkURLResponse
	for _, req := range requests {
		url := strings.TrimSpace(req.OriginalURL)

		shortURL, err := s.Shortener.Set(userID.(string), url)
		if err != nil {
			shortURL, err = s.Shortener.GetExistURL(url, err)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "Ошибка при сокращении URL",
					"code":    http.StatusInternalServerError,
				})
				return
			}
			httpStatus = http.StatusConflict
		}

		responses = append(responses, BulkURLResponse{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	c.JSON(httpStatus, responses)
}

// Обработчик редиректа на оригинальный URL по короткому ID
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	shortID := c.Param("id")

	originalURL, err := s.Shortener.Get(shortID)
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			c.Status(http.StatusGone)
			return
		}
		c.String(http.StatusTemporaryRedirect, err.Error())
		return
	}

	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, originalURL)
}

// Логика пинга для проверки работы сервиса
func (s *RestAPI) Ping(ctx *gin.Context) {
	if err := s.Shortener.Ping(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Ошибка при проверке состояния сервиса"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Сервис работает"})
}

// Обработчик для получения всех URL пользователя
func (s *RestAPI) UserURLsHandler(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Пользователь не авторизован",
		})
		return
	}

	urls, err := s.Shortener.GetFullRep(userID.(string))
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			ctx.Status(http.StatusGone)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка при получении списка URL пользователя",
		})
		return
	}

	if len(urls) == 0 {
		ctx.Status(http.StatusNoContent)
		return
	}

	ctx.JSON(http.StatusOK, urls)
}

// Обработчик удаления URL пользователя
func (s *RestAPI) DeleteUserUrls(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Пользователь не авторизован",
		})
		return
	}

	var shortURLs []string
	if err := ctx.BindJSON(&shortURLs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Некорректный запрос",
			"error":   err.Error(),
		})
		return
	}

	if err := s.Shortener.DeleteURLsRep(userID.(string), shortURLs); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Ошибка при удалении URL",
			"error":   err.Error(),
		})
		return
	}

	ctx.Status(http.StatusAccepted)
}
