package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"

	"google.golang.org/grpc"
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

	// Инициализация контекста
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создание и настройка gRPC-сервера
	grpcServer := grpc.NewServer()
	defer grpcServer.GracefulStop()

	// Регистрируем сервисы на gRPC-сервере
	// examplepb.RegisterExampleServiceServer(grpcServer, exampleServiceInstance)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Не удалось запустить gRPC listener: %v", err)
	}

	go func() {
		log.Println("Запуск gRPC-сервера на порту 50051...")
		if err := grpcServer.Serve(listener); err != nil {
			log.Fatalf("Ошибка запуска gRPC-сервера: %v", err)
		}
	}()

	// Обработка системных сигналов
	go func() {
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		<-signalChan
		log.Println("Получен сигнал завершения, инициирован shutdown...")
		cancel()
	}()

	// Запуск основного приложения
	if err := appInstance.Start(ctx); err != nil {
		log.Fatalf("Ошибка при запуске приложения: %v", err)
	}

	// Ожидание завершения контекста
	select {
	case <-ctx.Done():
		log.Println("Контекст отменен, завершение работы...")
	}
}
