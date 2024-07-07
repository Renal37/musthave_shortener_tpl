package config

import (
	"flag"
)

// Для запуска  ./shortener -a localhost:8888 -b http://localhost:8888

// Config хранит конфигурацию для сервера.
type Config struct {
	Address string
	BaseURL string
}

// InitConfig инициализирует конфигурацию из флагов командной строки.
func InitConfig() *Config {
	address := flag.String("a", "localhost:8080", "Адрес запуска HTTP-сервера")
	baseURL := flag.String("b", "http://localhost:8080", "Базовый адрес результирующего сокращённого URL")

	flag.Parse()

	return &Config{
		Address: *address,
		BaseURL: *baseURL,
	}
}
