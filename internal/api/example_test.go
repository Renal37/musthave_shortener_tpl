package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	// "strings"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api" // Импортируем пакет с RestAPI
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// // Пример использования ShortenURLHandler для сокращения URL
// func ExampleRestAPI_ShortenURLHandler() {
// 	storageInstance := storage.NewStorage()
// 	shortenerService := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
// 	apiInstance := &api.RestAPI{ // Используем api.RestAPI
// 		StructService: shortenerService,
// 	}

// 	// Инициализация маршрутов
// 	router := gin.Default()
// 	router.POST("/", apiInstance.ShortenURLHandler)

// 	// Создание запроса на сокращение URL
// 	w := httptest.NewRecorder()
// 	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("https://practicum.yandex.ru/"))

// 	// Отправка запроса
// 	router.ServeHTTP(w, req)

// 	// Проверка результата
// 	assert.Equal(nil, http.StatusCreated, w.Code)
// 	assert.Equal(nil, "text/plain", w.Header().Get("Content-Type"))
// 	assert.True(nil, strings.HasPrefix(w.Body.String(), "http://localhost:8080/"))
// }

// Пример использования ShortenURLJSON для сокращения URL с телом запроса в формате JSON
func ExampleRestAPI_ShortenURLJSON() {
	storageInstance := storage.NewStorage()
	shortenerService := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	apiInstance := &api.RestAPI{ // Используем api.RestAPI
		StructService: shortenerService,
	}

	// Инициализация маршрутов
	router := gin.Default()
	router.POST("/api/shorten", apiInstance.ShortenURLJSON)

	// Формирование тела запроса в формате JSON
	body := map[string]string{"url": "https://practicum.yandex.ru"}
	bodyJSON, _ := json.Marshal(body)

	// Создание HTTP-запроса
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	router.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(nil, http.StatusCreated, w.Code)
	assert.Equal(nil, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(nil, w.Body.String(), `"result"`)
}

// Пример использования ShortenURLsJSON для сокращения нескольких URL-адресов
func ExampleRestAPI_ShortenURLsJSON() {
	storageInstance := storage.NewStorage()
	shortenerService := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	apiInstance := &api.RestAPI{ // Используем api.RestAPI
		StructService: shortenerService,
	}

	// Инициализация маршрутов
	router := gin.Default()
	router.POST("/api/shorten", apiInstance.ShortenURLsJSON)

	// Формирование тела запроса с несколькими URL
	body := []map[string]string{
		{"correlation_id": "1", "original_url": "https://google.com"},
		{"correlation_id": "2", "original_url": "https://practicum.yandex.ru"},
	}
	bodyJSON, _ := json.Marshal(body)

	// Создание HTTP-запроса
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/shorten", bytes.NewReader(bodyJSON))
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	router.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(nil, http.StatusCreated, w.Code)
	assert.Equal(nil, "application/json", w.Header().Get("Content-Type"))
	assert.Contains(nil, w.Body.String(), `"short_url"`)
}

// Пример использования RedirectToOriginalURL для перенаправления по сокращенному URL
func ExampleRestAPI_RedirectToOriginalURL() {
	storageInstance := storage.NewStorage()
	shortenerService := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	apiInstance := &api.RestAPI{ // Используем api.RestAPI
		StructService: shortenerService,
	}

	// Добавляем тестовые данные в хранилище
	shortenerService.Storage.Set("abc123", "https://practicum.yandex.ru")

	// Инициализация маршрутов
	router := gin.Default()
	router.GET("/:id", apiInstance.RedirectToOriginalURL)

	// Создание HTTP-запроса для редиректа
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/abc123", nil)

	// Отправка запроса
	router.ServeHTTP(w, req)

	// Проверка результата
	assert.Equal(nil, http.StatusTemporaryRedirect, w.Code)
	assert.Equal(nil, "https://practicum.yandex.ru", w.Header().Get("Location"))
}
