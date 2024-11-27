package main

import (
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
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

	// Инициализация конфигурации
	addrConfig := config.InitConfig()

	// Инициализация хранилища и приложения
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запуск pprof сервера, если включен флаг
	if addrConfig.EnablePprof == "true" {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Основной процесс приложения
	go func() {
		err := appInstance.Start()
		if err != nil {
			log.Fatalf("Application error: %v\n", err)
		}
	}()

	// Создаем новый экземпляр gin.Engine
	router := gin.Default()

	// Настраиваем маршруты
	apiInstance := &api.RestAPI{} // Создайте экземпляр вашего API
	apiInstance.SetRoutes(router)

	// Запуск сервера (HTTP или HTTPS)
	if addrConfig.EnableHTTPS {
		fmt.Printf("Starting server with HTTPS at %s\n", addrConfig.ServerAddr)
		err := http.ListenAndServeTLS(
			addrConfig.ServerAddr,
			addrConfig.CertFile,
			addrConfig.KeyFile,
			router,
		)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v\n", err)
		}
	} else {
		fmt.Printf("Starting server without HTTPS at %s\n", addrConfig.ServerAddr)
		err := http.ListenAndServe(addrConfig.ServerAddr, router)
		if err != nil {
			log.Fatalf("Failed to start HTTP server: %v\n", err)
		}
	}

	// Завершаем приложение
	defer appInstance.Stop()
}
