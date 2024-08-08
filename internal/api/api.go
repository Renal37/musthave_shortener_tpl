package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/store"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RestAPI структура содержит указатель на объект ShortenerService
type RestAPI struct {
	StructService *services.ShortenerService
}

// StartRestAPI функция инициализирует сервер REST API
// Она принимает три параметра: ServerAddr, BaseURL, и storage
// ServerAddr - адрес, на котором будет запущен сервер
// BaseURL - основной URL для работы с коротким URL-адресатором
// storage - объект хранилища для короткого URL-адресатора
func StartRestAPI(ServerAddr, BaseURL string, LogLevel string, DBPath string, storage *storage.Storage) error {

	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}
	logger.Log.Info("Запуск сервера", zap.String("address", ServerAddr))
	bd, err := store.InitDatabase(DBPath)
	if err != nil {
		return err
	}
	// Создаем новый объект ShortenerService с указанным BaseURL и хранилищем
	storageShortener := services.NewShortenerService(BaseURL, storage, bd)

	// Создаем новый объект RestAPI с указанным ShortenerService
	api := &RestAPI{
		StructService: storageShortener,
	}

	// Устанавливаем режим работы Gin-фреймворка на ReleaseMode
	gin.SetMode(gin.ReleaseMode)
	// Создаем новый экземпляр Gin-инженерии
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.LoggerMiddleware(logger.Log),
		middleware.CompressMiddleware(),
	)

	api.setRoutes(r)

	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		err := r.Run(ServerAddr)
		if err != nil {
			fmt.Println("failed to start the browser")
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка при остановке сервера: %v\n", err)
	}

	return nil
}
