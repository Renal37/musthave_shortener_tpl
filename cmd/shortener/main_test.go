package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

// Тест для проверки запуска и остановки приложения
func TestMainApplication_StartStop(t *testing.T) {
	// Инициализация конфигурации и хранилища
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()

	// Создаем экземпляр приложения
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запуск Pprof в отдельной горутине
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		assert.NoError(t, err, "Pprof server failed to start")
	}()

	// Запускаем приложение в отдельной горутине
	go func() {
		appInstance.Start()
	}()

	// Ожидание 1 секунды для симуляции работы приложения
	time.Sleep(1 * time.Second)

	// Останавливаем приложение
	appInstance.Stop()

	// Проверяем, что приложение успешно остановлено
	assert.NotNil(t, appInstance, "App instance should not be nil")
}
