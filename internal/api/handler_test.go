package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Test_shortenURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	type args struct {
		code        int
		contentType string
	}
	tests := []struct {
		name    string
		Storage RestAPI
		args    args
		body    string
	}{
		{
			name: "test1",
			Storage: RestAPI{
				Shortener: storageShortener,
			},
			args: args{
				code:        201,
				contentType: "text/plain",
			},
			body: "https://practicum.yandex.ru/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()
			r.POST("/", tt.Storage.ShortenURLHandler)
			request := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			r.ServeHTTP(w, request)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, tt.args.code, res.StatusCode)
			assert.Equal(t, tt.args.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_shortenURLHandlerJSON(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	type args struct {
		code        int
		contentType string
	}
	type reqBody struct {
		PerformanceURL string `json:"url"`
	}
	tests := []struct {
		name    string
		Storage RestAPI
		args    args
		body    reqBody
	}{
		{
			name: "test1",
			Storage: RestAPI{
				Shortener: storageShortener,
			},
			args: args{
				code:        201,
				contentType: "application/json",
			},
			body: reqBody{
				"https://practicum.yandex.ru",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()

			r.POST("/api/shorten", tt.Storage.ShortenURLJSON)
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

			assert.Equal(t, tt.args.code, res.StatusCode)
			assert.Equal(t, tt.args.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_shortenURLsHandlerJSON(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	type args struct {
		code        int
		contentType string
	}
	type RequestBodyURLs struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}
	tests := []struct {
		name    string
		Storage RestAPI
		args    args
		body    []RequestBodyURLs
	}{
		{
			name: "test1",
			Storage: RestAPI{
				Shortener: storageShortener,
			},
			args: args{
				code:        201,
				contentType: "application/json",
			},
			body: []RequestBodyURLs{
				{
					CorrelationID: "1",
					OriginalURL:   "google.com",
				},
				{
					CorrelationID: "2",
					OriginalURL:   "google.kz",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := gin.Default()

			r.POST("/api/shorten", tt.Storage.ShortenURLsJSON)
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

			assert.Equal(t, tt.args.code, res.StatusCode)
			assert.Equal(t, tt.args.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_redirectToOriginalURLHandler(t *testing.T) {
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)

	type argsGet struct {
		code     int
		testURL  string
		location string
	}
	testsGET := []struct {
		name    string
		Storage RestAPI
		argsGet argsGet
	}{
		{
			name: "test1",
			Storage: RestAPI{
				Shortener: storageShortener,
			},
			argsGet: argsGet{
				code:     307,
				testURL:  "ads",
				location: "https://practicum.yandex.ru/",
			},
		},
	}

	for _, tt := range testsGET {
		t.Run(tt.name, func(t *testing.T) {
			tt.Storage.Shortener.Storage.Set(tt.argsGet.testURL, tt.argsGet.location)

			r := gin.Default()
			r.GET("/:id", tt.Storage.RedirectToOriginalURL)

			request := httptest.NewRequest(http.MethodGet, "/ads", nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, request)
			res := w.Result()
			defer res.Body.Close()
			assert.Equal(t, tt.argsGet.code, res.StatusCode)
			assert.Equal(t, tt.argsGet.location, res.Header.Get("location"))
		})
	}
}
