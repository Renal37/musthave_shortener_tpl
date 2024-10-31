// example_test.go
package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
)

// ExampleStartRestAPI демонстрирует, как запустить REST API.
// Этот пример не запускает настоящий сервер, а создает тестовый сервер для демонстрации.
func ExampleStartRestAPI() {
	// Создание тестового хранилища и базы данных
	storage := &storage.Storage{}
	db := &repository.StoreDB{}

	// Запуск API на тестовом адресе
	err := StartRestAPI(":8080", "http://localhost:8080", "info", db, false, storage)
	if err != nil {
		fmt.Println("Ошибка запуска API:", err)
		return
	}

	// Создание тестового запроса
	req, _ := http.NewRequest("GET", "/some-endpoint", nil)
	w := httptest.NewRecorder()

	// Создание нового роутера для тестирования
	r := gin.Default()
	r.GET("/some-endpoint", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
	})

	// Выполнение запроса
	r.ServeHTTP(w, req)

	// Проверка результата
	fmt.Println(w.Body.String()) // Вывод: {"message":"Hello, world!"}

	// Output:
	// {"message":"Hello, world!"}
}
