package api

import (
    "fmt"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/services"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
    "github.com/gin-gonic/gin"
)

// RestAPI структура содержит указатель на объект ShortenerService
type RestAPI struct {
    StructService *services.ShortenerService
}

// StartRestAPI функция инициализирует сервер REST API
// Она принимает три параметра: ServerAddr, BaseURL, и storage
// ServerAddr - адрес, на котором будет запущен сервер
// BaseURL - основной URL для работы с коротким URL-адресатором
// storage - объект хранилища для короткого URL-адресатора
func StartRestAPI(ServerAddr, BaseURL string, storage *storage.Storage) error {
    // Создаем новый объект ShortenerService с указанным BaseURL и хранилищем
    storageShortener := services.NewShortenerService(BaseURL, storage)
    // Создаем новый объект RestAPI с указанным ShortenerService
    api := &RestAPI{
        StructService: storageShortener,
    }

    // Устанавливаем режим работы Gin-фреймворка на ReleaseMode
    gin.SetMode(gin.ReleaseMode)
    // Создаем новый экземпляр Gin-инженерии
    r := gin.Default()

    // Вызываем метод setRoutes на объекте RestAPI для добавления маршрутов в API
    api.setRoutes(r)

    // Запускаем сервер на указанном ServerAddr
    err := r.Run(ServerAddr)
    // Если возникнет ошибка при запуске сервера, выводим сообщение об ошибке и возвращаем эту ошибку
    if err != nil {
        fmt.Println("Ошибка при запуске сервера: ", err)
        return err
    }

    // Если сервер запустился без ошибок, возвращаем nil
    return nil
}