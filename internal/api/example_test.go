package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Test shortening URLs with a handler
func TestShortenURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := api.RestAPI{Shortener: storageShortener}

	tests := []struct {
		name         string
		body         string
		expectedCode int
	}{
		{"Valid URL", "https://practicum.yandex.ru/", 201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/", handler.ShortenURLHandler)

			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)

			// Check status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Check if response body is a valid URL
			assert.Regexp(t, `^http://localhost:8080/[0-9a-fA-F-]{36}$`, w.Body.String())
		})
	}
}

// Test shortening URLs with JSON
func TestShortenURLHandlerJSON(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := api.RestAPI{Shortener: storageShortener}

	tests := []struct {
		name         string
		body         map[string]string
		expectedCode int
	}{
		{"Valid JSON URL", map[string]string{"url": "https://practicum.yandex.ru"}, 201},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/api/shorten", handler.ShortenURLJSON)

			jsonBody, _ := json.Marshal(tt.body)
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(jsonBody))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			// Check status code
			assert.Equal(t, tt.expectedCode, w.Code)

			// Validate response JSON structure
			var result map[string]string
			if err := json.Unmarshal(w.Body.Bytes(), &result); err == nil {
				assert.Regexp(t, `^http://localhost:8080/[0-9a-fA-F-]{36}$`, result["result"])
			}
		})
	}
}

// Test redirecting to the original URL
func TestRedirectToOriginalURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := api.RestAPI{Shortener: storageShortener}

	tests := []struct {
		name       string
		shortID    string
		original   string
		statusCode int
	}{
		{"Redirect to original URL", "test_id", "https://practicum.yandex.ru/", 307},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storageInstance.Set(tt.shortID, tt.original)

			r := gin.Default()
			r.GET("/:id", handler.RedirectToOriginalURL)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortID, nil)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			// Check status code
			assert.Equal(t, tt.statusCode, w.Code)
			assert.Equal(t, tt.original, w.Header().Get("Location"))
		})
	}
}
