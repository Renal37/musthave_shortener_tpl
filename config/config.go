package config

import (
    "flag"
    "os"
)

// Config представляет конфигурацию приложения.
type Config struct {
    ServerAddress string
    BaseURL       string
}

// InitConfig инициализирует конфигурацию, используя переменные окружения, флаги командной строки и значения по умолчанию.
func InitConfig() *Config {
    // Определение флагов командной строки
    serverAddressFlag := flag.String("a", "localhost:8080", "HTTP server address")
    baseURLFlag := flag.String("b", "http://localhost:8080", "Base URL for shortened links")

    // Разбор флагов командной строки
    flag.Parse()

    // Считывание переменных окружения
    serverAddressEnv := os.Getenv("SERVER_ADDRESS")
    baseURLEnv := os.Getenv("BASE_URL")

    // Установка конфигурации с приоритетом переменных окружения
    cfg := &Config{
        ServerAddress: *serverAddressFlag,
        BaseURL:       *baseURLFlag,
    }

    if serverAddressEnv != "" {
        cfg.ServerAddress = serverAddressEnv
    }
    if baseURLEnv != "" {
        cfg.BaseURL = baseURLEnv
    }

    return cfg
}

// Для запуска
// export SERVER_ADDRESS="localhost:8081"
// export BASE_URL="/myshorten"
// go run main.gox  