package api

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenURLHandler обрабатывает запросы на сокращение URL без JSON
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
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "Не удалось сократить URL", http.StatusInternalServerError)
			return
		}
		httpStatus = http.StatusConflict
	}
	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

// ShortenURLJSON обрабатывает запросы на сокращение URL в формате JSON
func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
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

	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	url := strings.TrimSpace(decoderBody.URL)
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
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

	StructPerformance := Response{Result: shortURL}
	respJSON, err := json.Marshal(StructPerformance)
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Не удалось прочитать тело запроса",
			"code":    http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorMassage)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	code := http.StatusTemporaryRedirect
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
	c.String(code, originalURL)
}

// ShortenURLsJSON обрабатывает запросы на сокращение нескольких URL в формате JSON
func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
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

	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL)
		shortURL, err := s.Shortener.Set(userID, url)
		if err != nil {
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
		urlResponse := ResponseBodyURLs{
			req.CorrelationID,
			shortURL,
		}
		URLResponses = append(URLResponses, urlResponse)
	}
	respJSON, err := json.Marshal(URLResponses)
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
	c.Data(httpStatus, "application/json", respJSON)
}

// Ping проверяет доступность службы
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.Shortener.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}

// UserURLsHandler возвращает URL-адреса пользователя
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
		ctx.JSON(code, nil)
		return
	}
	userID, _ := userIDFromContext.(string)
	urls, err := s.Shortener.GetFullRep(userID)
	ctx.Header("Content-type", "application/json")
	if err != nil {
		if err.Error() == http.StatusText(http.StatusGone) {
			ctx.Status(http.StatusGone)
			return
		}
		code = http.StatusInternalServerError
		ctx.JSON(code, gin.H{
			"message": "Не удалось получить URL-адреса пользователя",
			"code":    code,
		})
		return
	}

	if len(urls) == 0 {
		ctx.JSON(http.StatusNoContent, nil)
		return
	}
	ctx.JSON(code, urls)
}

// DeleteUserUrls удаляет URL-адреса пользователя
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
	//проверка
	var shortURLs []string
	if err := ctx.BindJSON(&shortURLs); err != nil {
		code = http.StatusBadRequest
		ctx.JSON(code, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := s.Shortener.DeleteURLsRep(userID, shortURLs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось удалить URL-адрес",
			"error":   errors.New("не удалось удалить URL-адрес").Error(),
		})
		return
	}
	ctx.Status(code)
}
