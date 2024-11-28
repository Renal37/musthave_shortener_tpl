package main

import (
	"crypto/tls"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func main() {
	// Вывод информации о сборке
	printBuildInfo()

	// Инициализация конфигурации
	cfg := config.InitConfig()

	// Запуск pprof, если включен
	if cfg.EnablePprof == "true" {
		go startPprofServer()
	}

	// Создаем маршруты
	router := gin.Default()
	apiInstance := &api.RestAPI{}
	apiInstance.SetRoutes(router)

	// Запускаем сервер
	if cfg.EnableHTTPS {
		startHTTPSServer(cfg, router)
	} else {
		startHTTPServer(cfg, router)
	}
}

// printBuildInfo выводит информацию о версии сборки
func printBuildInfo() {
	buildVersion := getEnvOrDefault("BUILD_VERSION", "N/A")
	buildDate := getEnvOrDefault("BUILD_DATE", "N/A")
	buildCommit := getEnvOrDefault("BUILD_COMMIT", "N/A")

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// startPprofServer запускает pprof сервер на порту 6060
func startPprofServer() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

// startHTTPSServer запускает HTTPS сервер
func startHTTPSServer(cfg *config.Config, router http.Handler) {
	checkFileExists(cfg.CertFile, "SSL certificate")
	checkFileExists(cfg.KeyFile, "SSL key")

	tlsConfig := &tls.Config{MinVersion: tls.VersionTLS12}
	server := &http.Server{
		Addr:      cfg.ServerAddr,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	log.Printf("Starting HTTPS server at %s...\n", cfg.ServerAddr)
	if err := server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile); err != nil {
		log.Fatalf("Error starting HTTPS server: %v", err)
	}
}

// startHTTPServer запускает HTTP сервер
func startHTTPServer(cfg *config.Config, router http.Handler) {
	log.Printf("Starting HTTP server at %s...\n", cfg.ServerAddr)
	if err := http.ListenAndServe(cfg.ServerAddr, router); err != nil {
		log.Fatalf("Error starting HTTP server: %v", err)
	}
}

// checkFileExists проверяет существование файла
func checkFileExists(path string, description string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatalf("%s file %s does not exist", description, path)
	}
}

// getEnvOrDefault возвращает значение переменной окружения или значение по умолчанию
func getEnvOrDefault(env, defaultValue string) string {
	if value := os.Getenv(env); value != "" {
		return value
	}
	return defaultValue
}
