package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	LogLevel   string `env:"FLAG_LOG_LEVEL"`
	FilePath   string `env:"FILE_STORAGE_PATH"`
}

// Функция инициализации конфигурации
func InitConfig() *Config {
	config := &Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		LogLevel:   "info",
		FilePath:   "short-url-db.json",
	}
	// Установка флагов для адреса и порта сервера так же ошибка
	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "адрес и номер порта для запуска API")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "адрес и номер порта для запуска API адресПозиции")
	flag.StringVar(&config.LogLevel, "c", config.LogLevel, "log level")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "address to file in-memory")

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
