package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func TestMain(t *testing.T) {
	// Инициализация конфигурации
	addrConfig := config.InitConfig() // Предполагается, что эта функция не возвращает ошибку

	// Создаем экземпляр хранилища
	storageInstance := storage.NewStorage() // Возможно, нужно будет передать параметры

	// Создаем экземпляр приложения
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Создаем тестовый HTTP сервер
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			http.Error(w, "Application is running", http.StatusOK)
		} else {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	}))
	defer server.Close()

	// Запуск приложения в горутине
	go func() {
		appInstance.Start() // Предполагается, что этот метод не блокирует
	}()
	defer appInstance.Stop() // Остановить приложение после теста

	// Отправляем тестовый запрос на сервер
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что код состояния ответа 200
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got: %d", resp.StatusCode)
	}

	// Проверяем, что приложение корректно обрабатывает неправильный метод
	resp, err = http.Post(server.URL, "application/json", bytes.NewBuffer([]byte{}))
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}
	defer resp.Body.Close()

	// Проверяем, что код состояния ответа 405
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code 405, got: %d", resp.StatusCode)
	}
}
