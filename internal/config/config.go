package config

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"sync"
)

// Config представляет собой структуру для хранения конфигурационных параметров приложения.
type Config struct {
	ServerAddr  string `env:"SERVER_ADDRESS"`    // Адрес и порт сервера
	BaseURL     string `env:"BASE_URL"`          // Базовый URL
	LogLevel    string `env:"FLAG_LOG_LEVEL"`    // Уровень логирования
	FilePath    string `env:"FILE_STORAGE_PATH"` // Путь к файлу для хранения
	DBPath      string `env:"DB_PATH"`           // Путь к базе данных
	EnablePprof string `env:"ENABLE_PPROF"`      // Включить pprof
	EnableHTTPS string `env:"ENABLE_HTTPS"`      // Включить HTTPS
	CertFile    string `env:"CERT_FILE"`         // Путь к файлу сертификата
	KeyFile     string `env:"KEY_FILE"`          // Путь к файлу ключа
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
		EnablePprof: "false",                 // Значение по умолчанию для pprof
		EnableHTTPS: "false",                 // Значение по умолчанию для включения HTTPS
		CertFile:    "server.crt",            // Значение по умолчанию для сертификата
		KeyFile:     "server.key",            // Значение по умолчанию для ключа
	}

	// Определяем флаги командной строки
	once.Do(func() {
		flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "address and port to run api")
		flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base URL")
		flag.StringVar(&config.LogLevel, "c", config.LogLevel, "log level")
		flag.StringVar(&config.FilePath, "f", config.FilePath, "path to file for storage")
		flag.StringVar(&config.DBPath, "d", config.DBPath, "path to database")
		flag.StringVar(&config.EnablePprof, "e", config.EnablePprof, "enable pprof")
		flag.StringVar(&config.EnableHTTPS, "s", config.EnableHTTPS, "enable https (true/false)")
		flag.StringVar(&config.CertFile, "cert", config.CertFile, "path to the SSL certificate file")
		flag.StringVar(&config.KeyFile, "key", config.KeyFile, "path to the SSL key file")
		flag.Parse() // Парсим флаги командной строки
	})

	// Загружаем параметры из переменных окружения
	err := env.Parse(config)
	if err != nil {
		panic(err) // Завершаем программу, если произошла ошибка
	}

	return config // Возвращаем структуру конфигурации
}
