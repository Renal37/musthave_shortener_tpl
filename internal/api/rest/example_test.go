package rest

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
)

// Example для обработки сокращения URL
func ShortenURLHandlers() {
	// Инициализация зависимостей
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := RestAPI{Shortener: storageShortener}

	// Настройка роутера Gin
	r := gin.Default()
	r.POST("/", handler.ShortenURLHandler)

	// Создание HTTP-запроса
	body := "https://practicum.yandex.ru/"
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	// Выполнение запроса
	r.ServeHTTP(w, request)

	// Выводим результат
	fmt.Println(w.Code)
	fmt.Println(w.Body.String())

	// Output:
	// 201
	// http://localhost:8080/<UUID>
}

// Example для редиректа на оригинальный URL
func RedirectToOriginalURLHandlers() {
	// Инициализация зависимостей
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := RestAPI{Shortener: storageShortener}
	// Добавляем URL в хранилище
	storageInstance.Set("test_id", "https://practicum.yandex.ru/")

	// Настройка роутера Gin
	r := gin.Default()
	r.GET("/:id", handler.RedirectToOriginalURL)

	// Создание HTTP-запроса
	request := httptest.NewRequest(http.MethodGet, "/test_id", nil)
	w := httptest.NewRecorder()

	// Выполнение запроса
	r.ServeHTTP(w, request)

	// Выводим результат
	fmt.Println(w.Code)
	fmt.Println(w.Header().Get("Location"))

	// Output:
	// 307
	// https://practicum.yandex.ru/
}
