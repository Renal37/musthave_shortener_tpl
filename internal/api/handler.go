package api

import (
	"encoding/json"
	"errors"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/gin-gonic/gin"
	"io"
	"net"
	"net/http"
	"strings"
)

// Request представляет структуру для обработки запроса на сокращение URL
type Request struct {
	URL string `json:"url"`
}

// Response представляет структуру для ответа с сокращенным URL
type Response struct {
	Result string `json:"result"`
}

// RequestBodyURLs представляет запрос с уникальным идентификатором корреляции
// и оригинальным URL для обработки сокращения
type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// ResponseBodyURLs представляет ответ с уникальным идентификатором корреляции
// и сокращенным URL
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// ShortenURLHandler обрабатывает запросы на сокращение URL, переданного в теле запроса в виде строки.
// Возвращает сокращенный URL. Если URL уже существует, возвращает имеющийся сокращенный URL
// со статусом 409 Conflict.
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

// ShortenURLJSON обрабатывает запросы на сокращение URL в формате JSON.
// Возвращает JSON с сокращенным URL. Если URL уже существует, возвращает имеющийся сокращенный URL
// со статусом 409 Conflict.
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
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
		shortURL, err = s.Shortener.GetExistURL(url, err)
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
			"message": "Не удалось создать ответ",
			"code":    http.StatusInternalServerError,
		}
		c.JSON(http.StatusInternalServerError, errorMessage)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// RedirectToOriginalURL перенаправляет пользователя на оригинальный URL по сокращенному идентификатору.
// Если URL недоступен, возвращает статус 410 Gone.
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

// ShortenURLsJSON обрабатывает запросы на сокращение нескольких URL в формате JSON.
// Возвращает JSON со списком сокращенных URL с их идентификаторами корреляции.
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
		answer, _ := json.Marshal(errorMessage)
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
		var answer []byte
		answer, err = json.Marshal(errorMessage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Внутренняя ошибка сервера"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Ping проверяет доступность службы, возвращая статус 200 OK в случае успешного ответа.
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.Shortener.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}

// UserURLsHandler возвращает все URL-адреса, созданные пользователем.
// Если пользователь не найден, возвращает статус 401 Unauthorized.
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

// DeleteUserUrls удаляет URL-адреса пользователя на основе списка сокращенных URL-адресов.
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

// StatsHandler обрабатывает запросы на получение статистики о сервисе сокращения URL.
// Проверяет, что запрос исходит из доверенной подсети, и возвращает общее количество
// сокращённых URL и количество пользователей. Ответ предоставляется в формате JSON.
func (s *RestAPI) StatsHandler(c *gin.Context) {
	addrConfig := config.InitConfig()
	trustedSubnet := addrConfig.TrustedSubnet
	clientIP := c.GetHeader("X-Real-IP")

	if trustedSubnet != "" {
		_, cidr, err := net.ParseCIDR(trustedSubnet)
		if err != nil || !cidr.Contains(net.ParseIP(clientIP)) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			return
		}
	}

	urlCount, err := s.Shortener.GetURLCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch URL count"})
		return
	}

	userCount, err := s.Shortener.GetUserCount()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch user count"})
		return
	}

	stats := map[string]int{
		"urls":  urlCount,
		"users": userCount,
	}

	c.JSON(http.StatusOK, stats)
}
