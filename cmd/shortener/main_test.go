package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock components for testing
type MockStorage struct {
	mock.Mock
}

type MockService struct {
	mock.Mock
}

func TestInitializeAndStartApp(t *testing.T) {
	// Mock configuration
	mockConfig := &config.Config{
		ServerAddr: ":8080",
		BaseURL:    "http://localhost:8080",
		FilePath:   "/tmp/test_data.json",
	}

	// Mock storage
	mockStorage := &storage.Storage{}

	// Mock services
	mockService := services.NewShortenerService(mockConfig.BaseURL, mockStorage, nil, false)

	// Создаем экземпляр приложения
	appInstance := app.NewApp(mockStorage, mockService, mockConfig)

	// Настраиваем контекст и канал для сигналов
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signalChan := make(chan os.Signal, 1)
	defer close(signalChan)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		// Эмулируем сигнал завершения через небольшую задержку
		time.Sleep(1 * time.Second)
		signalChan <- syscall.SIGINT
	}()

	// Запускаем приложение в отдельной горутине
	errChan := make(chan error, 1)
	go func() {
		errChan <- appInstance.Start(ctx)
	}()

	// Ожидаем сигнал завершения
	select {
	case <-signalChan:
		cancel()
	case <-time.After(2 * time.Second):
		t.Fatal("Тест завершился по таймауту при ожидании сигнала")
	}

	// Проверяем, что приложение завершилось без ошибок
	select {
	case err := <-errChan:
		assert.NoError(t, err, "Приложение должно завершиться без ошибок")
	case <-time.After(1 * time.Second):
		t.Fatal("Тест завершился по таймауту при ожидании завершения приложения")
	}
}
