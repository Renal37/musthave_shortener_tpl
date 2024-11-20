package main

import (
	"context"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	buildVersion string // Версия сборки
	buildDate    string // Дата сборки
	buildCommit  string // Хэш коммита
)

func main() {
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

	// Инициализация конфигурации, хранилища и приложения
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Создаем маршрутизатор Gin
	router := gin.Default()

	// Инициализация маршрутов через приложение

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:    addrConfig.ServerAddr,
		Handler: router,
	}

	// Запуск pprof сервера на порту 6060
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Запуск основного HTTP-сервера
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске приложения: %v", err)
		}
	}()

	// Создаем канал для получения сигналов
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Блокируем выполнение до получения сигнала
	<-quit

	// Останавливаем приложение
	appInstance.Stop()

	// Завершаем работу сервера с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при завершении работы сервера: %v", err)
	}

	log.Println("Сервер завершил работу")
}
