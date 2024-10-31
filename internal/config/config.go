package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// Config представляет собой структуру для хранения конфигурационных параметров приложения.
type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`    // Адрес и порт сервера
	BaseURL    string `env:"BASE_URL"`          // Базовый URL
	LogLevel   string `env:"FLAG_LOG_LEVEL"`    // Уровень логирования
	FilePath   string `env:"FILE_STORAGE_PATH"` // Путь к файлу для хранения
	DBPath     string `env:"db"`                // Путь к базе данных
}

// InitConfig инициализирует конфигурацию, загружая параметры из переменных окружения и флагов командной строки.
func InitConfig() *Config {
	config := &Config{
		ServerAddr: "localhost:8080",        // Значение по умолчанию для адреса сервера
		BaseURL:    "http://localhost:8080", // Значение по умолчанию для базового URL
		LogLevel:   "info",                  // Значение по умолчанию для уровня логирования
		FilePath:   "short-url-db.json",     // Значение по умолчанию для пути к файлу
		DBPath:     "",                      // Значение по умолчанию для пути к базе данных
	}

	// Определяем флаги командной строки
	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "address and port to run api")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "address and port to run api addrResPos")
	flag.StringVar(&config.LogLevel, "c", config.LogLevel, "log level")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "address to file in-memory")
	flag.StringVar(&config.DBPath, "d", config.DBPath, "address to base store in-memory")

	flag.Parse() // Парсим флаги командной строки

	// Загружаем параметры из переменных окружения
	err := env.Parse(config)
	if err != nil {
		panic(err) // Завершаем программу, если произошла ошибка
	}

	return config // Возвращаем структуру конфигурации
}
