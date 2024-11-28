package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
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
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запуск pprof сервера на порту 6060, если включен флаг
	if addrConfig.EnablePprof == "true" {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}
	r := gin.Default()
	// Определяем запуск сервера (HTTP или HTTPS)
	if addrConfig.EnableHTTPS {
		// Запуск HTTPS-сервера
		fmt.Println("Starting server in HTTPS mode...")
		err := http.ListenAndServeTLS(addrConfig.ServerAddr, addrConfig.CertFile, addrConfig.KeyFile, r)
		if err != nil {
			log.Fatalf("Failed to start HTTPS server: %v", err)
		}
	} else {
		// Запуск HTTP-сервера
		fmt.Println("Starting server in HTTP mode...")
		appInstance.Start()
	}

	// Завершение работы приложения
	appInstance.Stop()
}
