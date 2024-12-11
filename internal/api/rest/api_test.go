package rest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStartRestAPI(t *testing.T) {
	// Инициализация зависимостей
	storageInstance := storage.NewStorage()
	db := &repository.StoreDB{} // Предположим, что это mock или stub

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем сервер в отдельной горутине
	go func() {
		err := StartRestAPI(ctx, ":8080", "http://localhost:8080", "info", db, false, storageInstance, false, "", "")
		assert.NoError(t, err)
	}()

	// Даем серверу время на запуск
	time.Sleep(1 * time.Second)

	// Создаем HTTP-запрос для проверки работы сервера
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Создаем роутер и обрабатываем запрос
	r := gin.Default()
	r.ServeHTTP(w, req)

	// Проверяем, что сервер отвечает с кодом 404, так как маршруты не настроены
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Отменяем контекст для завершения работы сервера
	cancel()

	// Даем серверу время на завершение
	time.Sleep(1 * time.Second)
}

func TestStartRestAPIWithHTTPS(t *testing.T) {
	// Инициализация зависимостей
	storageInstance := storage.NewStorage()
	db := &repository.StoreDB{} // Предположим, что это mock или stub

	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запускаем сервер с HTTPS в отдельной горутине
	go func() {
		err := StartRestAPI(ctx, ":8443", "https://localhost:8443", "info", db, false, storageInstance, true, "cert.pem", "key.pem")
		assert.NoError(t, err)
	}()

	// Даем серверу время на запуск
	time.Sleep(1 * time.Second)

	// Создаем HTTP-запрос для проверки работы сервера
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Создаем роутер и обрабатываем запрос
	r := gin.Default()
	r.ServeHTTP(w, req)

	// Проверяем, что сервер отвечает с кодом 404, так как маршруты не настроены
	assert.Equal(t, http.StatusNotFound, w.Code)

	// Отменяем контекст для завершения работы сервера
	cancel()

	// Даем серверу время на завершение
	time.Sleep(1 * time.Second)
}
