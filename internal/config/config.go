package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
)

// Config представляет собой структуру для хранения конфигурационных параметров приложения.
type Config struct {
	ServerAddr  string `env:"SERVER_ADDRESS" json:"server_address"`       // Адрес и порт сервера
	BaseURL     string `env:"BASE_URL" json:"base_url"`                   // Базовый URL
	LogLevel    string `env:"FLAG_LOG_LEVEL" json:"-"`                    // Уровень логирования (только флаг или env)
	FilePath    string `env:"FILE_STORAGE_PATH" json:"file_storage_path"` // Путь к файлу для хранения
	DBPath      string `env:"DB_PATH" json:"database_dsn"`                // Путь к базе данных
	EnablePprof string `env:"ENABLE_PPROF" json:"-"`                      // Включить pprof (только флаг или env)
	EnableHTTPS bool   `env:"ENABLE_HTTPS" json:"enable_https"`           // Включить HTTPS
	CertFile    string `env:"CERT_FILE" json:"cert_file"`                 // Путь к файлу сертификата
	KeyFile     string `env:"KEY_FILE" json:"key_file"`                   // Путь к файлу ключа
	ConfigPath  string `env:"CONFIG" json:"-"`  
	TrustedSubnet string `env:"TRUSTED_SUBNET" json:"trusted_subnet"`                          // Путь к файлу конфигурации (только флаг или env)
}

var once sync.Once

// LoadConfigFromFile загружает конфигурацию из файла JSON, если файл указан.
func LoadConfigFromFile(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// mergeConfigs объединяет данные из JSON-конфигурации и основной конфигурации (где данные из флагов и env имеют приоритет).
func mergeConfigs(base, fileConfig *Config) *Config {
	if fileConfig.ServerAddr != "" {
		base.ServerAddr = fileConfig.ServerAddr
	}
	if fileConfig.BaseURL != "" {
		base.BaseURL = fileConfig.BaseURL
	}
	if fileConfig.FilePath != "" {
		base.FilePath = fileConfig.FilePath
	}
	if fileConfig.DBPath != "" {
		base.DBPath = fileConfig.DBPath
	}
	if fileConfig.EnablePprof != "" {
		base.EnablePprof = fileConfig.EnablePprof
	}
	if fileConfig.EnableHTTPS {
		base.EnableHTTPS = fileConfig.EnableHTTPS
	}
	if fileConfig.CertFile != "" {
		base.CertFile = fileConfig.CertFile
	}
	if fileConfig.KeyFile != "" {
		base.KeyFile = fileConfig.KeyFile
	}
	return base
}

// InitConfig инициализирует конфигурацию, загружая параметры из файла, переменных окружения и флагов командной строки.
func InitConfig() *Config {
	config := &Config{
		ServerAddr:  "localhost:8080",        // Значение по умолчанию для адреса сервера
		BaseURL:     "http://localhost:8080", // Значение по умолчанию для базового URL
		LogLevel:    "info",                  // Значение по умолчанию для уровня логирования
		FilePath:    "short-url-db.json",     // Значение по умолчанию для пути к файлу
		DBPath:      "",                      // Значение по умолчанию для пути к базе данных
		EnablePprof: "false",                 // Значение по умолчанию для pprof
		EnableHTTPS: false,                   // Значение по умолчанию для HTTPS
		CertFile:    "cert.pem",              // Значение по умолчанию для сертификата
		KeyFile:     "key.pem",               // Значение по умолчанию для ключа
	}

	// Определяем флаги командной строки
	once.Do(func() {
		flag.StringVar(&config.ServerAddr, "a", config.ServerAddr, "address and port to run api")
		flag.StringVar(&config.BaseURL, "b", config.BaseURL, "base URL")
		flag.StringVar(&config.LogLevel, "c", config.LogLevel, "log level")
		flag.StringVar(&config.FilePath, "f", config.FilePath, "path to file for storage")
		flag.StringVar(&config.DBPath, "d", config.DBPath, "path to database")
		flag.StringVar(&config.EnablePprof, "e", config.EnablePprof, "enable pprof")
		flag.BoolVar(&config.EnableHTTPS, "s", config.EnableHTTPS, "enable https (true/false)")
		flag.StringVar(&config.CertFile, "cert", config.CertFile, "path to the SSL certificate file")
		flag.StringVar(&config.KeyFile, "key", config.KeyFile, "path to the SSL key file")
		flag.StringVar(&config.ConfigPath, "config", config.ConfigPath, "path to config file")
		flag.StringVar(&config.TrustedSubnet, "t", config.TrustedSubnet, "доверенная подсеть в формате CIDR")
		flag.Parse() // Парсим флаги командной строки
	})

	// Загружаем параметры из переменных окружения
	err := env.Parse(config)
	if err != nil {
		panic(fmt.Sprintf("Error parsing env variables: %v", err))
	}

	// Если указан файл конфигурации, загружаем его
	if config.ConfigPath != "" {
		fileConfig, err := LoadConfigFromFile(config.ConfigPath)
		if err != nil {
			fmt.Printf("Error loading config file: %v\n", err)
		} else {
			// Объединяем конфигурацию из файла с основной (данные из файла имеют меньший приоритет)
			config = mergeConfigs(config, fileConfig)
		}
	}

	return config // Возвращаем структуру конфигурации
}
