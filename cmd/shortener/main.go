package main

import (
	"context"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"log"
	_ "net/http/pprof"
)

var (
	buildVersion string // Версия сборки
	buildDate    string // Дата сборки
	buildCommit  string // Хэш коммита
)

func main() {
	initializeAndStartApp()
}

func initializeAndStartApp() {
	// Если переменные не были переданы при компиляции, выводим "N/A"
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	// Выводим информацию о сборке
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// Инициализация конфигурации
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запуск приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := appInstance.Start(ctx); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
