package api

import (
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/services"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
    "net/http"
    "net/http/httptest"
    "testing"
)

// Test_shortenURLHandler - тест для функции ShortenURLHandler
func Test_shortenURLHandler(t *testing.T) {
    // Описание теста
    type args struct {
        code        int
        contentType string
    }
    tests := []struct {
        name    string
        Storage RestAPI
        args    args
    }{
        {
            name: "test1",
            Storage: RestAPI{
                StructService: &services.ShortenerService{
                    Storage: &storage.Storage{},
                },
            },
            args: args{
                code:        201,
                contentType: "text/plain",
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Очищаем URLs перед каждым тестом
            tt.Storage.StructService.Storage.URLs = make(map[string]string)
            tt.Storage.StructService.BaseURL = "http://localhost:8081"

            // Создаем объект gin.Engine и устанавливаем маршрут для теста
            r := gin.Default()

            // Вызываем функцию ShortenURLHandler и проверяем результат
            r.POST("/", tt.Storage.ShortenURLHandler)

            request := httptest.NewRequest(http.MethodPost, "/", nil)
            w := httptest.NewRecorder()

            r.ServeHTTP(w, request)

            res := w.Result()
            defer res.Body.Close()

            // Проверяем статус ответа и тип контента
            assert.Equal(t, tt.args.code, res.StatusCode)
            assert.Equal(t, tt.args.contentType, res.Header.Get("Content-Type"))
        })
}
}

// Test_redirectToOriginalURLHandler - тест для функции RedirectToOriginalURLHandler
func Test_redirectToOriginalURLHandler(t *testing.T) {
    // Описание теста
    type argsGet struct {
        code     int
        testURL  string
        location string
    }
    testsGET := []struct {
        name    string
        Storage RestAPI
        argsGet argsGet
    }{
        {
            name: "test1",
            Storage: RestAPI{
                    StructService: &services.ShortenerService{
                        Storage: &storage.Storage{},
                    },
            },
            argsGet: argsGet{
                    code:     307,
                    testURL:  "ads",
                    location: "https://practicum.yandex.ru/",
                },
        },
    }

    for _, tt := range testsGET {
        t.Run(tt.name, func(t *testing.T) {
            // Устанавливаем базовый URL для тестов
            tt.Storage.StructService.BaseURL = "http://localhost:8081"

            // Очищаем URLs перед каждым тестом
            tt.Storage.StructService.Storage.URLs = make(map[string]string)
            tt.Storage.StructService.Storage.Set(tt.argsGet.testURL, tt.argsGet.location)

            // Создаем объект gin.Engine и устанавливаем маршрут для теста
            r := gin.Default()
            r.GET("/:id", tt.Storage.RedirectToOriginalURLHandler)

            // Проверяем, что функция возвращает правильный статус и локацию
            request := httptest.NewRequest(http.MethodGet, "/ads", nil)
            w := httptest.NewRecorder()
            r.ServeHTTP(w, request)
            res := w.Result()
            defer res.Body.Close()

            // Проверяем, что функция возвращает правильный статус и локацию
            assert.Equal(t, tt.argsGet.code, res.StatusCode)
            assert.Equal(t, tt.argsGet.location, res.Header.Get("location"))
        })
    }
}