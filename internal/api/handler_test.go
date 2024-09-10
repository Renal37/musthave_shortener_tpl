package api

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestShortenURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	tests := []struct {
		name        string
		body        string
		expectedCode int
		expectedContentType string
	}{
		{
			name: "valid URL",
			body: "https://practicum.yandex.ru/",
			expectedCode: http.StatusCreated,
			expectedContentType: "text/plain",
		},
		{
			name: "empty body",
			body: "",
			expectedCode: http.StatusInternalServerError,
			expectedContentType: "text/plain",
		},
		{
			name: "invalid URL",
			body: "not a valid url",
			expectedCode: http.StatusInternalServerError,
			expectedContentType: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			api := RestAPI{Shortener: storageShortener}
			r.POST("/", api.ShortenURLHandler)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedContentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestShortenURLJSON(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	tests := []struct {
		name        string
		body        URLRequest
		expectedCode int
		expectedContentType string
	}{
		{
			name: "valid URL",
			body: URLRequest{URL: "https://practicum.yandex.ru"},
			expectedCode: http.StatusCreated,
			expectedContentType: "application/json",
		},
		{
			name: "empty URL",
			body: URLRequest{URL: ""},
			expectedCode: http.StatusInternalServerError,
			expectedContentType: "application/json",
		},
		{
			name: "invalid URL",
			body: URLRequest{URL: "not a valid url"},
			expectedCode: http.StatusInternalServerError,
			expectedContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			api := RestAPI{Shortener: storageShortener}
			r.POST("/api/shorten", api.ShortenURLJSON)
			jsonBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "/api/shorten", strings.NewReader(string(jsonBody)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedContentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestShortenURLsHandlerJSON(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	tests := []struct {
		name        string
		body        []BulkURLRequest
		expectedCode int
		expectedContentType string
	}{
		{
			name: "valid URLs",
			body: []BulkURLRequest{
				{CorrelationID: "1", OriginalURL: "https://google.com"},
				{CorrelationID: "2", OriginalURL: "https://google.kz"},
			},
			expectedCode: http.StatusCreated,
			expectedContentType: "application/json",
		},
		{
			name: "empty URL",
			body: []BulkURLRequest{
				{CorrelationID: "1", OriginalURL: ""},
			},
			expectedCode: http.StatusInternalServerError,
			expectedContentType: "application/json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			api := RestAPI{Shortener: storageShortener}
			r.POST("/api/shorten/bulk", api.ShortenURLsJSON)
			jsonBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			request := httptest.NewRequest(http.MethodPost, "/api/shorten/bulk", strings.NewReader(string(jsonBody)))
			request.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			assert.Equal(t, tt.expectedContentType, res.Header.Get("Content-Type"))
		})
	}
}

func TestRedirectToOriginalURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	tests := []struct {
		name        string
		shortID     string
		originalURL string
		expectedCode int
	}{
		{
			name: "valid shortID",
			shortID: "short-id",
			originalURL: "https://practicum.yandex.ru/",
			expectedCode: http.StatusTemporaryRedirect,
		},
		{
			name: "nonexistent shortID",
			shortID: "nonexistent-id",
			expectedCode: http.StatusTemporaryRedirect,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.originalURL != "" {
				storageShortener.Storage.Set(tt.shortID, tt.originalURL)
			}

			r := gin.Default()
			api := RestAPI{Shortener: storageShortener}
			r.GET("/:id", api.RedirectToOriginalURL)

			request := httptest.NewRequest(http.MethodGet, "/"+tt.shortID, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.expectedCode, res.StatusCode)
			if tt.originalURL != "" {
				assert.Equal(t, tt.originalURL, res.Header.Get("Location"))
			}
		})
	}
}
