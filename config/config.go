package config

import (
    "os"
)

type Config struct {
    Address     string
    BaseURL     string
}
// Для запуска
// export SERVER_ADDRESS="localhost:8081"
// export BASE_URL="/myshorten"
// go run main.go


func InitConfig() Config {
    cfg := Config{
        Address:     os.Getenv("SERVER_ADDRESS"),
        BaseURL:     os.Getenv("BASE_URL"),
    }

    if cfg.Address == "" {
        cfg.Address = "localhost:8080"
    }

    if cfg.BaseURL == "" {
        cfg.BaseURL = "/shorten"
    }

    return cfg
}