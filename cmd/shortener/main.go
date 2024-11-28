package main

import (
	"crypto/tls"
	"flag"
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
	// Если переменные не были переданы при компиляции, выводим "N/A"
	buildVersion := "N/A"
	buildDate := "N/A"
	buildCommit := "N/A"

	// Выводим информацию о сборке
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	// Обработка флага -enableHTTPS
	enableHTTPS := flag.Bool("enableHTTPS", false, "Enable HTTPS")
	flag.Parse()

	// Проверка переменной окружения ENABLE_HTTPS
	if envHTTPS := os.Getenv("ENABLE_HTTPS"); envHTTPS == "true" {
		*enableHTTPS = true
	}

	addrConfig := config.InitConfig() // Инициализация конфигурации

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
	apiInstance := &api.RestAPI{}
	apiInstance.SetRoutes(r)

	// Проверка наличия файлов сертификата и ключа
	certFile := "server.crt"
	keyFile := "server.key"

	if *enableHTTPS {
		if _, err := os.Stat(certFile); os.IsNotExist(err) {
			log.Fatalf("Certificate file %s does not exist", certFile)
		}

		if _, err := os.Stat(keyFile); os.IsNotExist(err) {
			log.Fatalf("Key file %s does not exist", keyFile)
		}

		// Создание HTTPS сервера
		server := &http.Server{
			Addr:      addrConfig.ServerAddr,
			Handler:   r,
			TLSConfig: tlsConfig,
		}

		// Запуск HTTPS сервера
		err := server.ListenAndServeTLS(certFile, keyFile)
		if err != nil {
			log.Fatalf("Error starting HTTPS server: %v", err)
		}
	} else {
		// Запуск HTTP сервера
		err := http.ListenAndServe(addrConfig.ServerAddr, r)
		if err != nil {
			log.Fatalf("Error starting HTTP server: %v", err)
		}
	}
}