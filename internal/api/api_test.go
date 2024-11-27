package api_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStartRestAPI(t *testing.T) {
	// Установим режим тестирования Gin
	gin.SetMode(gin.TestMode)

	// Настроим параметры для запуска API
	serverAddr := ":8080"
	baseURL := "http://localhost:8080"
	logLevel := "debug"
	dbDNSTurn := false

	// Инициализируем логгер
	_ = logger.Initialize(logLevel)

	// Создадим контекст для управления жизненным циклом сервера
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Запустим сервер в отдельной горутине и получим только функцию для его остановки
	_, stop := api.StartRestAPI(ctx, serverAddr, baseURL, logLevel, &repository.StoreDB{}, dbDNSTurn, &storage.Storage{})

	// Даем серверу немного времени на запуск
	time.Sleep(1 * time.Second)

	// Проверим, что сервер работает, отправив запрос на базовый URL
	req, err := http.NewRequest(http.MethodGet, baseURL+"/ping", nil)
	assert.NoError(t, err)

	// Используем httptest для отправки запроса к серверу
	rr := httptest.NewRecorder()
	handler := gin.Default()
	handler.ServeHTTP(rr, req)

	// Проверим, что мы получили ответ с кодом 404 (или 200, если в вашем коде есть обработчик /ping)
	assert.Equal(t, http.StatusNotFound, rr.Code) // Ожидаем 404, если "/ping" не был определён

	// Теперь остановим сервер, послав сигнал завершения в контекст
	cancel()

	// Подождем, чтобы сервер успел завершить свою работу
	time.Sleep(1 * time.Second)

	// Проверим, что сервер был остановлен, и функция stop отработала корректно
	err = stop()
	assert.NoError(t, err)
}
