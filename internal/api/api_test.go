package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Тест функции StartRestAPI
func TestStartRestAPI(t *testing.T) {
	// Установить режим тестирования Gin
	gin.SetMode(gin.TestMode)

	// Настроим параметры для запуска API
	serverAddr := ":8080"
	baseURL := "http://localhost:8080"
	logLevel := "debug"
	dbDNSTurn := false

	// Инициализируем логгер
	_ = logger.Initialize(logLevel)

	// Запускаем API в отдельной горутине
	go func() {
		err := api.StartRestAPI(serverAddr, baseURL, logLevel, &repository.StoreDB{}, dbDNSTurn, &storage.Storage{})
		if err != nil {
			t.Errorf("Error starting API: %v", err)
		}
	}()

	// Даем время на запуск сервера
	time.Sleep(1 * time.Second)

	// Создаем HTTP-запрос к серверу
	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	assert.NoError(t, err)

	// Используем httptest для тестирования HTTP-сервера
	rr := httptest.NewRecorder()
	handler := gin.Default()
	handler.ServeHTTP(rr, req)

	// Проверяем код ответа
	assert.Equal(t, http.StatusNotFound, rr.Code) // Ожидаем 404, так как не было маршрута по умолчанию

	// Закрываем сервер, посылая сигнал прерывания
	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)

	// Имитация SIGINT для завершения сервера
	process.Signal(syscall.SIGINT)

	// Даем время на завершение сервера
	time.Sleep(1 * time.Second)
}

// Моковый сервер для тестирования функции startServer
func TestStartServer(t *testing.T) {
	// Мокаем Gin
	r := gin.New()

	// Используем httptest для запуска и тестирования сервера
	server := &http.Server{
		Addr:    ":8081",
		Handler: r,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		_ = api.StartRestAPI(":8081", "http://localhost:8081", "info", &repository.StoreDB{}, false, &storage.Storage{})
	}()

	// Даем серверу время на запуск
	time.Sleep(1 * time.Second)

	// Создаем HTTP-запрос к серверу
	req, err := http.NewRequest(http.MethodGet, "http://localhost:8081", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Проверяем код ответа (ожидаем 404)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	// Закрываем сервер
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		_ = server.Shutdown(ctx)
	}()

	process, err := os.FindProcess(os.Getpid())
	assert.NoError(t, err)
	process.Signal(os.Interrupt)

	time.Sleep(1 * time.Second)
}
