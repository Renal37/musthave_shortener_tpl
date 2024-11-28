package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
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

	addrConfig := config.InitConfig()                      // Инициализация конфигурации
	storageInstance := storage.NewStorage()                // Создание хранилища
	appInstance := app.NewApp(storageInstance, addrConfig) // Создание приложения

	// Запуск pprof сервера на порту 6060, если включен флаг
	if addrConfig.EnablePprof == "true" {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	// Настройка TLS конфигурации
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Создание экземпляра gin.Engine
	r := gin.Default()

	// Настройка маршрутов
	apiInstance := &api.RestAPI{} // Создайте экземпляр RestAPI
	apiInstance.SetRoutes(r)

	// Проверка наличия файлов сертификата и ключа
	certFile := "server.crt"
	keyFile := "server.key"

	if _, err := os.Stat(certFile); os.IsNotExist(err) {
		log.Fatalf("Certificate file %s does not exist", certFile)
	}

	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		log.Fatalf("Key file %s does not exist", keyFile)
	}

	// Создание HTTPS сервера
	server := &http.Server{
		Addr:      addrConfig.ServerAddr,
		Handler:   r, // Используем gin.Engine в качестве обработчика
		TLSConfig: tlsConfig,
	}

	// Запуск HTTPS сервера
	err := server.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatalf("Error starting HTTPS server: %v", err)
	}

	appInstance.Start()
	appInstance.Stop()
}
