package config

// Импортируем необходимые пакеты
import (
	"flag"
	"github.com/caarlos0/env/v6"
)

// Определяем структуру Config, которая будет хранить параметры конфигурации
type Config struct {
	ServerAddr string `env:"SERVER_ADDRESS"`
	BaseURL    string `env:"BASE_URL"`
	LogLevel   string `env:"FLAG_LOG_LEVEL"`
	FilePath   string `env:"FILE_STORAGE_PATH"`
	DATABASE_DSN     string `env:"db"`
}

// Функция InitConfig инициализирует структуру Config с помощью флагов и переменных окружения
func InitConfig() *Config {
	config := &Config{
		ServerAddr: "localhost:8080",
		BaseURL:    "http://localhost:8080",
		LogLevel:   "info",
		FilePath:   "short-url-db.json",
		DATABASE_DSN:     "",
	}

	// Устанавливаем флаги для параметров конфигурации
	flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "address and port to run api")
	flag.StringVar(&config.BaseURL, "b", config.BaseURL, "address and port to run api addrResPos")
	flag.StringVar(&config.LogLevel, "c", config.LogLevel, "log level")
	flag.StringVar(&config.FilePath, "f", config.FilePath, "address to file in-memory")
	flag.StringVar(&config.DATABASE_DSN, "d", config.DATABASE_DSN, "address to base store in-memory")

	// Парсим флаги и переменные окружения
	flag.Parse()
	err := env.Parse(config)
	if err != nil {
		panic(err)
	}

	// Возвращаем инициализированную структуру Config
	return config
}
