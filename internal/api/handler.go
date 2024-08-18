package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// Определение структуры Request для URL-адреса
type Request struct {
	URL string `json:"url"`
}

// Определение структуры Response для результата
type Response struct {
	Result string `json:"result"`
}

// Определение структуры RequestBodyURLs для корреляционного идентификатора и исходного URL-адреса
type RequestBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

// Определение структуры ResponseBodyURLs для корреляционного идентификатора и сокращенного URL-адреса
type ResponseBodyURLs struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

// Функция ShortenURLHandler обрабатывает запрос на сокращение URL-адреса
func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса", http.StatusInternalServerError)
		return
	}
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)
	s.StructService.UserID = userID
	url := strings.TrimSpace(string(body))
	shortURL, err := s.StructService.Set(url)
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

// Функция ShortenURLJSON обрабатывает запрос на сокращение URL-адреса в формате JSON
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
	s.StructService.UserID = userID
	url := strings.TrimSpace(decoderBody.URL)
	shortURL, err := s.StructService.Set(url)
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err)
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
		answer, _ := json.Marshal(errorMassage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Функция RedirectToOriginalURL обрабатывает запрос на перенаправление по сокращенному URL-адресу
func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	shortID := c.Param("id")
	originalURL, exists := s.StructService.Get(shortID)
	if !exists {
		c.String(http.StatusTemporaryRedirect, "URL не найден", http.StatusTemporaryRedirect)
		return
	}
	c.Header("Location", originalURL)
	c.String(http.StatusTemporaryRedirect, originalURL)
}

// Функция ShortenURLsJSON обрабатывает запрос на сокращение нескольких URL-адресов в формате JSON
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
		answer, _ := json.Marshal(errorMassage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL)
		shortURL, err := s.StructService.Set(url)
		if err != nil {
			shortURL, err = s.StructService.GetExistURL(url, err)
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
		answer, _ := json.Marshal(errorMassage)
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)
}

// Функция Ping обрабатывает запрос на проверку работоспособности сервиса
func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.StructService.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}
