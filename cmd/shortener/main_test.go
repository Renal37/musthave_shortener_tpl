package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLHandler(t *testing.T) {
	// Инициализация зависимостей
	storageInstance := storage.NewStorage()
	storageShortener := services.NewShortenerService("http://localhost:8080", storageInstance, nil, false)
	handler := api.RestAPI{Shortener: storageShortener}

	// Настройка роутера Gin
	r := gin.Default()
	r.POST("/", handler.ShortenURLHandler)

	// Создание HTTP-запроса
	body := "https://practicum.yandex.ru/"
	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	// Выполнение запроса
	r.ServeHTTP(w, request)

	// Проверка результатов
	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "http://localhost:8080/")
}
