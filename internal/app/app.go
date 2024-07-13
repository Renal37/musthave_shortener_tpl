package app

import (
    "fmt"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/api"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/config"
    "github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

// Функция Start - запускает сервер приложения
func Start(config *config.Config) {
    // Создаем экземпляр хранилища
    storageInstance := storage.NewStorage()

    // Запускаем сервер приложения
    err := api.StartRestAPI(config.ServerAddr, config.BaseURL, storageInstance)
    if err != nil {
        // Если возникла ошибка, выводим сообщение об ошибке и возвращаем
        fmt.Println("Ошибка при запуске сервера: ", err)
        return
    }
}