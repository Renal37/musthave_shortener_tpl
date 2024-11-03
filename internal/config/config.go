package config

import (
	"flag"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config представляет собой структуру для хранения конфигурационных параметров приложения.
type Config struct {
	ServerAddr  string `env:"SERVER_ADDRESS"`    // Адрес и порт сервера
	BaseURL     string `env:"BASE_URL"`          // Базовый URL
	LogLevel    string `env:"LOG_LEVEL"`         // Уровень логирования
	FilePath    string `env:"FILE_STORAGE_PATH"` // Путь к файлу для хранения
	DBPath      string `env:"DB_PATH"`           // Путь к базе данных
	EnablePprof bool   `env:"ENABLE_PPROF"`      // Включение профайлинга
}

var once sync.Once

// InitConfig инициализирует конфигурацию, загружая параметры из переменных окружения и флагов командной строки.
func InitConfig() *Config {
	config := &Config{
		ServerAddr:  "localhost:8080",        // Значение по умолчанию для адреса сервера
		BaseURL:     "http://localhost:8080", // Значение по умолчанию для базового URL
		LogLevel:    "info",                  // Значение по умолчанию для уровня логирования
		FilePath:    "short-url-db.json",     // Значение по умолчанию для пути к файлу
		DBPath:      "",                      // Значение по умолчанию для пути к базе данных
		EnablePprof: false,                   // Значение по умолчанию для включения профайлинга
	}

	// Определяем флаги командной строки
	once.Do(func() {
		flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "address and port to run api")
		flag.StringVar(&config.BaseURL, "-b", config.BaseURL, "base URL")
		flag.StringVar(&config.LogLevel, "-c", config.LogLevel, "log level")
		flag.StringVar(&config.FilePath, "-f", config.FilePath, "path to file for storage")
		flag.StringVar(&config.DBPath, "-d", config.DBPath, "path to database")
		flag.BoolVar(&config.EnablePprof, "e", config.EnablePprof, "enable profiling")
		flag.Parse() // Парсим флаги командной строки
	})

	// Загружаем параметры из переменных окружения
	if err := env.Parse(config); err != nil {
		panic(err) // Завершаем программу, если произошла ошибка
	}

	return config // Возвращаем структуру конфигурации
}
