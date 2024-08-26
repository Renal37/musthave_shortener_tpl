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

func (s *RestAPI) ShortenURLHandler(c *gin.Context) {
	httpStatus := http.StatusCreated
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	userIDFromContext, _ := c.Get("userID")
	userID, _ := userIDFromContext.(string)

	url := strings.TrimSpace(string(body))
	shortURL, err := s.Shortener.Set(userID, url)
	if err != nil {
		shortURL, err = s.Shortener.GetExistURL(url, err)
		if err != nil {
			c.String(http.StatusInternalServerError, "the url could not be shortened", http.StatusInternalServerError)
			return
		}
		httpStatus = http.StatusConflict
	}
	c.Header("Content-Type", "text/plain")
	c.String(httpStatus, shortURL)
}

func (s *RestAPI) ShortenURLJSON(c *gin.Context) {
	var decoderBody Request
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Failed to read request body",
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
				"message": "the url could not be shortened",
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
			"message": "Failed to read request body",
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

func (s *RestAPI) ShortenURLsJSON(c *gin.Context) {
	var decoderBody []RequestBodyURLs
	httpStatus := http.StatusCreated
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&decoderBody)
	c.Header("Content-Type", "application/json")
	if err != nil {
		errorMassage := map[string]interface{}{
			"message": "Failed to read request body",
			"code":    http.StatusInternalServerError,
		}
		var answer []byte
		answer, err = json.Marshal(errorMassage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
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
					"message": "the url could not be shortened",
					"code":    http.StatusInternalServerError,
				}
				var answer []byte
				answer, err = json.Marshal(errorMassage)
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
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
			"message": "Failed to read request body",
			"code":    http.StatusInternalServerError,
		}
		var answer []byte
		answer, err = json.Marshal(errorMassage)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
			return
		}
		c.Data(http.StatusInternalServerError, "application/json", answer)
		return
	}
	c.Data(httpStatus, "application/json", respJSON)

}

func (s *RestAPI) Ping(ctx *gin.Context) {
	err := s.Shortener.Ping()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "")
		return
	}
	ctx.JSON(http.StatusOK, "")
}

func (s *RestAPI) UserURLsHandler(ctx *gin.Context) {
	code := http.StatusOK
	userIDFromContext, exists := ctx.Get("userID")
	ctx.Request.Context()
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get userID",
			"error":   errors.New("failed to get user from context").Error(),
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
			"message": "Failed to retrieve user URLs",
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

func (s *RestAPI) DeleteUserUrls(ctx *gin.Context) {
	code := http.StatusAccepted
	userIDFromContext, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get userID",
			"error":   errors.New("failed to get user from context").Error(),
		})
		return
	}
	userID, _ := userIDFromContext.(string)

	var shorURLs []string
	if err := ctx.BindJSON(&shorURLs); err != nil {
		code = http.StatusBadRequest
		ctx.JSON(code, gin.H{
			"error:": err.Error(),
		})
	}

	err := s.Shortener.DeleteURLsRep(userID, shorURLs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed delete to url",
			"error":   errors.New("failed to get user from context").Error(),
		})
		return
	}
	ctx.Status(code)
}