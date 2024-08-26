package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

// Структура запроса для обработки URL.
type Request struct {
	URL string `json:"url"`
}

// Структура ответа для отправки сокращенного URL.
type Response struct {
	Result string `json:"result"`
}

// Структура тела запроса для обработки нескольких URL.
type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Структура тела ответа для нескольких URL.
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenURLHandler обрабатывает запрос на сокращение URL.
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	url := strings.TrimSpace(string(body))
	shortURL, err := s.StructService.Set(userID, url)
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "Не удалось сократить URL", http.StatusInternalServerError)
			return
		}
		httpStatus = http.StatusConflict
	}
	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

// ShortenURLJSON обрабатывает JSON-запрос на сокращение URL.
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
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

	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	url := strings.TrimSpace(decoderBody.URL)
	shortURL, err := s.StructService.Set(userID, url)
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err)
		if err != nil {
			errorMessage := map[string]interface{}{
				"message": "Не удалось сократить URL",
				"code":    http.StatusInternalServerError,
			}
			answer, _ := json.Marshal(errorMessage)
			c.Data(http.StatusInternalServerError, "application/json", answer)
			return
		}
		httpStatus = http.StatusConflict
	}

	response := Response{Result: shortURL}
	respJSON, err := json.Marshal(response)
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// RedirectToOriginalURL обрабатывает перенаправление на оригинальный URL.
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	code := http.StatusTemporaryRedirect
	shortID := c.Param("id")
	originalURL, err := s.StructService.Get(shortID)
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			c.Status(http.StatusGone)
			return
		}
		c.String(http.StatusTemporaryRedirect, err.Error())
		return
	}

	c.Header("Location", originalURL)
	c.String(code, originalURL)
}

// ShortenURLsJSON обрабатывает JSON-запрос на сокращение нескольких URL.
func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, err := json.Marshal(errorMessage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}

	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL)
		shortURL, err := s.StructService.Set(userID, url)
		if err != nil {
			shortURL, err = s.StructService.GetExistURL(url, err)
			if err != nil {
				errorMessage := map[string]interface{}{
					"message": "Не удалось сократить URL",
					"code":    http.StatusInternalServerError,
				}
				answer, err := json.Marshal(errorMessage)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
					return
				}
				c.Data(http.StatusInternalServerError, "application/json", answer)
				return
			}
			httpStatus = http.StatusConflict
		}
		urlResponse := ResponseBodyURLs{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		}
		URLResponses = append(URLResponses, urlResponse)
	}
	respJSON, err := json.Marshal(URLResponses)
	if err != nil {
		errorMessage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, err := json.Marshal(errorMessage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Ping проверяет доступность сервиса.
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.StructService.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Сервис доступен"})
}

// UserURLsHandler возвращает все URL пользователя.
func (s *RestAPI) UserURLsHandler(ctx *gin.Context) {
	code := http.StatusOK
	userIDFromContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось получить userID",
			"error":   errors.New("не удалось получить пользователя из контекста").Error(),
		})
		return
	}
	UserNew, _ := ctx.Get("new")
	if UserNew == true {
		code = http.StatusUnauthorized
		ctx.JSON(code, gin.H{"message": "Пользователь не авторизован"})
		return
	}
	userID, _ := userIDFromContext.(string)
	urls, err := s.StructService.GetFullRep(userID)
	ctx.Header("Content-type", "application/json")
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			ctx.Status(http.StatusGone)
			return
		}
		code = http.StatusInternalServerError
		ctx.JSON(code, gin.H{
			"message": "Не удалось получить URL пользователя",
			"code":    code,
		})
		return
	}

	if len(urls) == 0 {
		ctx.JSON(http.StatusNoContent, gin.H{"message": "Нет содержимого"})
		return
	}
	ctx.JSON(code, urls)
}

// DeleteUserUrls помечает ссылки как удаленные для определенного пользователя.
func (s *RestAPI) DeleteUserUrls(ctx *gin.Context) {
	code := http.StatusAccepted
	userIDFromContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось получить userID",
			"error":   errors.New("не удалось получить пользователя из контекста").Error(),
		})
		return
	}
	userID, _ := userIDFromContext.(string)

	var shortURLs []string
	if err := ctx.BindJSON(&shortURLs); err != nil {
		code = http.StatusBadRequest
		ctx.JSON(code, gin.H{
			"message": "Неверный запрос",
			"error:":  err.Error(),
		})
		return
	}

	err := s.StructService.DeleteURLsRep(userID, shortURLs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось удалить URL",
			"error":   err.Error(),
		})
		return
	}
	ctx.Status(code)
}
