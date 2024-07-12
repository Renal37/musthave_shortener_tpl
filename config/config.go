package config

import (
	"flag"
	"fmt"
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
	serverAddressFlag := flag.String("a", "", "HTTP server address")
	baseURLFlag := flag.String("b", "", "Base URL for shortened links")

	// Разбор флагов командной строки
	flag.Parse()

	// Считывание переменных окружения
	serverAddressEnv := os.Getenv("SERVER_ADDRESS")
	baseURLEnv := os.Getenv("BASE_URL")

	// Установка конфигурации с приоритетом флагов командной строки
	cfg := &Config{
		ServerAddress: "localhost:8080",
		BaseURL:       "http://localhost:8080",
	}

	if serverAddressEnv != "" {
		cfg.ServerAddress = serverAddressEnv
	}
	if baseURLEnv != "" {
		cfg.BaseURL = baseURLEnv
	}
	if *serverAddressFlag != "" {
		cfg.ServerAddress = *serverAddressFlag
	}
	if *baseURLFlag != "" {
		cfg.BaseURL = *baseURLFlag
	}
	// Проверка на nil перед использованием
	if cfg == nil {
		fmt.Println("Config is nil. Please check your code.")
		os.Exit(1)
	}
	return cfg
}

// Для запуска
// export SERVER_ADDRESS="localhost:8081"
// export BASE_URL="http://localhost:8081"
// go run main.go -a "localhost:8082" -b "http://localhost:8082"
