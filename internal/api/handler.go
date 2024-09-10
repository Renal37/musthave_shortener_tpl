package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Структура для получения URL из тела запроса
type Request struct {
	URL string `json:"url"`
}

// Структура для ответа с коротким URL
type Response struct {
	Result string `json:"result"`
}

// Структура для тела запроса с множеством URL
type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Структура для ответа с множеством коротких URL
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	shortURL      string `json:"short_url"`
}

// Обработчик сокращения URL из текстового тела запроса
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated

	// Читаем тело запроса
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		// Логируем и возвращаем ошибку, если не удалось прочитать тело
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}

	// Получаем userID из контекста
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	// Удаляем лишние пробелы из URL
	url := strings.TrimSpace(string(body))

	// Пытаемся сохранить URL
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
		// Если произошла ошибка, пытаемся получить существующий короткий URL
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "Не удалось сократить URL", http.StatusInternalServerError)
			return
		}
		// Если URL уже существует, устанавливаем статус 409 (Conflict)
		httpStatus = http.StatusConflict
	}

	// Устанавливаем заголовок и возвращаем короткий URL
	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

// Обработчик сокращения URL в формате JSON
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)

	// Парсим тело запроса
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		answer, _ := json.Marshal(errorMassage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}

	// Получаем userID из контекста
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	url := strings.TrimSpace(decoderBody.URL)

	// Пытаемся сохранить URL
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
		// Если произошла ошибка, пытаемся получить существующий короткий URL
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			errorMassage := map[string]interface{}{
				"message": "Не удалось сократить URL",
				"code":    http.StatusInternalServerError,
			}
			answer, _ := json.Marshal(errorMassage)
			c.Data(http.StatusInternalServerError, "application/json", answer)
			return
		}
		httpStatus = http.StatusConflict
	}

	// Создаем ответную структуру
	StructPerformance := Response{Result: shortURL}
	respJSON, err := json.Marshal(StructPerformance)
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Не удалось создать ответ",
			"code":    http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorMassage)
		return
	}

	// Отправляем JSON с коротким URL
	c.Data(httpStatus, "application/json", respJSON)
}

// Обработчик редиректа на оригинальный URL по короткому ID
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	code := http.StatusTemporaryRedirect
	shortID := c.Param("id")

	// Пытаемся получить оригинальный URL
	originalURL, err := s.Shortener.Get(shortID)
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			// Если URL был удален, возвращаем 410 (Gone)
			c.Status(http.StatusGone)
			return
		}
		// Если произошла другая ошибка, возвращаем её текст
		c.String(http.StatusTemporaryRedirect, err.Error())
		return
	}

	// Устанавливаем заголовок Location и делаем редирект
	c.Header("Location", originalURL)
	c.String(code, originalURL)
}

// Обработчик сокращения множества URL (JSON формат)
func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)

	// Парсим массив URL из тела запроса
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		var answer []byte
		answer, err = json.Marshal(errorMassage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}

	// Получаем userID из контекста
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	// Формируем ответы для каждого URL
	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL)

		// Пытаемся сократить URL
		shortURL, err := s.Shortener.Set(userID, url)
		if err != nil {
			// Если ошибка, проверяем, существует ли уже короткий URL
			shortURL, err = s.Shortener.GetExistURL(url, err)
			if err != nil {
				errorMassage := map[string]interface{}{
					"message": "Не удалось сократить URL",
					"code":    http.StatusInternalServerError,
				}
				var answer []byte
				answer, err = json.Marshal(errorMassage)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
					return
				}
				c.Data(http.StatusInternalServerError, "application/json", answer)
				return
			}
			httpStatus = http.StatusConflict
		}

		// Добавляем ответ в массив
		urlResponse := ResponseBodyURLs{
			CorrelationID: req.CorrelationID,
			shortURL:      shortURL,
		}
		URLResponses = append(URLResponses, urlResponse)
	}

	// Отправляем JSON с результатами
	respJSON, err := json.Marshal(URLResponses)
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Не удалось создать ответ",
			"code":    http.StatusInternalServerError,
		}
		var answer []byte
		answer, err = json.Marshal(errorMassage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Логика пинга для проверки работы сервиса
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.Shortener.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}

// Обработчик для получения всех URL пользователя
func (s *RestAPI) UserURLsHandler(ctx *gin.Context) {
	// Проверяем, что userID есть в контексте
	userIDFromContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Пользователь не авторизован",
		})
		return
	}

	// Получаем userID
	userID, _ := userIDFromContext.(string)

	// Получаем список всех URL пользователя
	urls, err := s.Shortener.GetFullRep(userID)
	ctx.Header("Content-type", "application/json")
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			ctx.Status(http.StatusGone)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось получить список URL пользователя",
		})
		return
	}

	// Если список пустой, возвращаем 204 (No Content)
	if len(urls) == 0 {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}

	// Возвращаем список URL
	ctx.JSON(http.StatusOK, urls)
}

// Обработчик удаления URL пользователя
func (s *RestAPI) DeleteUserUrls(ctx *gin.Context) {
	// Проверяем, что userID есть в контексте
	userIDFromContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Пользователь не авторизован",
		})
		return
	}

	// Получаем userID
	userID, _ := userIDFromContext.(string)

	// Читаем список коротких URL для удаления
	var shortURLs []string
	if err := ctx.BindJSON(&shortURLs); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Некорректный запрос",
			"error":   err.Error(),
		})
		return
	}

	// Пытаемся удалить URL
	err := s.Shortener.DeleteURLsRep(userID, shortURLs)
	if err != nil {
		// Логируем и возвращаем ошибку
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось удалить URL",
			"error":   err.Error(),
		})
		return
	}

	// Возвращаем успешный статус 202 (Accepted)
	ctx.Status(http.StatusAccepted)
}