package api

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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

func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Не удалось прочитать тело запроса")
		return
	}
	url := strings.TrimSpace(string(body))

	shortURL, err := s.StructService.Set(url)
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "Не удалось сократить URL")
			return
		}
		httpStatus = http.StatusConflict
	}

	c.Header("Content-Type", "application/json")
	c.String(httpStatus, shortURL)
}

func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось прочитать тело запроса", "code": http.StatusInternalServerError})
		return
	}
	url := strings.TrimSpace(decoderBody.URL)

	shortURL, err := s.StructService.Set(url)
	if err != nil {
		shortURL, err = s.StructService.GetExistURL(url, err)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось сократить URL", "code": http.StatusInternalServerError})
			return
		}
		httpStatus = http.StatusConflict
	}

	c.JSON(httpStatus, Response{Result: shortURL})
}

func (s *RestAPI) RedirectToOriginalURL(c *gin.Context) {
	shortID := c.Param("id")
	originalURL, exists := s.StructService.Get(shortID)
	if !exists {
		c.String(http.StatusNotFound, "URL не найден")
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, originalURL)
}

func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось прочитать тело запроса", "code": http.StatusInternalServerError})
		return
	}

	var URLResponses []ResponseBodyURLs
	for _, req := range decoderBody {
		url := strings.TrimSpace(req.OriginalURL)
		shortURL, err := s.StructService.Set(url)
		if err != nil {
			shortURL, err = s.StructService.GetExistURL(url, err)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Не удалось сократить URL", "code": http.StatusInternalServerError})
				return
			}
			httpStatus = http.StatusConflict
		}
		URLResponses = append(URLResponses, ResponseBodyURLs{
			CorrelationID: req.CorrelationID,
			ShortURL:      shortURL,
		})
	}

	c.JSON(httpStatus, URLResponses)
}

func (s *RestAPI) Ping(ctx *gin.Context) {
	if err := s.StructService.Ping(); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return
	}
	ctx.Status(http.StatusOK)
}
