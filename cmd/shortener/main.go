package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
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
	if buildVersion == "" {
		buildVersion = "N/A"
	}
	if buildDate == "" {
		buildDate = "N/A"
	}
	if buildCommit == "" {
		buildCommit = "N/A"
	}

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	// в InitDatabase("написать свой пусть к бд")
	dbInstance, err := repository.InitDatabase("")
	if err != nil {
		log.Fatalf("Ошибка при инициализации базы данных: %v", err)
	}

	dbDNSTurn := false // Установите в true, если хотите использовать базу данных для хранения

	servicesInstance := services.NewShortenerService("http://localhost:8080", storageInstance, dbInstance, dbDNSTurn)
	appInstance := app.NewApp(storageInstance, servicesInstance, addrConfig)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Канал для системных сигналов
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

		// Ожидание сигнала завершения
		<-signalChan
		cancel()
	}()

	if err := appInstance.Start(ctx); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}
}
