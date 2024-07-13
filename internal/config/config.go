package config

import (
    "flag"
    "fmt"
    "github.com/caarlos0/env/v6"
)

type Config struct {
    ServerAddr string `env:"SERVER_ADDRESS"`
    BaseURL    string `env:"BASE_URL"`
}

// Функция инициализации конфигурации
func InitConfig() *Config {
    config := &Config{
        ServerAddr: "localhost:8080",
        BaseURL:    "http://localhost:8080",
    }
    // Установка флагов для адреса и порта сервера
    flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "адрес и номер порта для запуска API")
    flag.StringVar(&config.BaseURL, "b", config.BaseURL, "адрес и номер порта для запуска API адресПозиции")
    // Проверка и парсинг переменных среды
    flag.Parse()
    err := env.Parse(config)
    if err != nil {
        // Если произошла ошибка при парсинге переменных среды, выводим сообщение об ошибке и возвращаем nil
        fmt.Printf("Ошибка при парсинге переменных среды: %v\n", err)
        return nil
    }
    // Возвращение инициализированной конфигурации
    return config
}