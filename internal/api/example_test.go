package api

import (
    "fmt"
    "net/http"
    "net/http/httptest"

    "github.com/gin-gonic/gin"
)

// MockStorage представляет собой пустую реализацию интерфейса storage.Storage,
// необходимую для запуска теста.
type MockStorage struct{}

// MockStoreDB представляет собой пустую реализацию интерфейса repository.StoreDB.
type MockStoreDB struct{}

// ExampleStartRestAPI демонстрирует, как можно запустить сервер REST API и выполнить к нему запрос.
// Этот пример не запускает постоянный сервер, а создает временный тестовый сервер.
func ExampleStartRestAPI() {
    // Создаем заглушки для storage и db
    storage := &MockStorage{}
    db := &MockStoreDB{}

    // Запуск API на тестовом адресе (можно изменить на нужный адрес)
    go func() {
        err := StartRestAPI(":8080", "http://example.com", "info", db, false, storage)
        if err != nil {
            fmt.Println("Ошибка запуска API:", err)
        }
    }()

    // Создание тестового запроса
    req, _ := http.NewRequest("GET", "/some-endpoint", nil)
    w := httptest.NewRecorder()

    // Настройка роутера для тестирования
    r := gin.Default()
    r.GET("/some-endpoint", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Hello, world!"})
    })

    // Выполнение тестового запроса
    r.ServeHTTP(w, req)

    // Выводим результат
    fmt.Println(w.Body.String())

    // Output:
    // {"message":"Hello, world!"}
}
