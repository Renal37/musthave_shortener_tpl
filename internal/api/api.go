package api

import (
	"context"
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// RestAPI представляет собой структуру для REST API.
type RestAPI struct {
	Shortener *services.ShortenerService // Сервис для сокращения URL.
}

// StartRestAPI запускает HTTP-сервер REST API для обработки запросов сокращения URL.
//
// Сервер создает и настраивает маршруты с использованием middleware и предоставляет эндпоинты
// для работы с короткими ссылками, а также обеспечивает автоматическое логирование запросов и управление авторизацией.
//
// Параметры:
//   - ServerAddr: Адрес для запуска сервера, например, ":8080".
//   - BaseURL: Базовый URL сервиса сокращения ссылок.
//   - LogLevel: Уровень логирования, например, "info" или "debug".
//   - db: Объект StoreDB для хранения данных в базе (предполагается реализация интерфейса базы данных).
//   - dbDNSTurn: Логический флаг, указывающий, использовать ли DNS-трансляцию для базы данных.
//   - storage: Объект Storage для хранения данных (предполагается реализация интерфейса хранилища).
//
// Пример использования:
//
//	go func() {
//	    err := StartRestAPI(":8080", "http://example.com", "info", db, false, storage)
//	    if err != nil {
//	        log.Fatalf("Ошибка запуска API: %v", err)
//	    }
//	}()
//
// Известные ограничения:
//   - Если сервер запущен, его завершение произойдет при получении сигнала остановки от системы.
//   - Время завершения при остановке ограничено 5 секундами (по умолчанию).
//
// BUG(Автор): Текущее логирование ограничено уровнем info; более детализированные уровни требуют дальнейшей настройки.
// BUG(Автор): В конфигурации сервера может отсутствовать поддержка HTTPS.

func StartRestAPI(ctx context.Context, ServerAddr, BaseURL, LogLevel string, db *repository.StoreDB, dbDNSTurn bool, storage *storage.Storage) error {
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}

	logger.Log.Info("Running server", zap.String("address", ServerAddr))
	storageShortener := services.NewShortenerService(BaseURL, storage, db, dbDNSTurn)

	api := &RestAPI{
		Shortener: storageShortener,
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(
		gin.Recovery(),
		middleware.LoggerMiddleware(logger.Log),
		middleware.CompressMiddleware(),
		middleware.AuthorizationMiddleware(),
	)

	api.SetRoutes(r)

	// Create an HTTP server
	srv := &http.Server{
		Addr:    ServerAddr,
		Handler: r,
	}

	// Start the server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the context cancellation and shutdown the server
	<-ctx.Done()

	logger.Log.Info("Shutting down server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	logger.Log.Info("Server exited properly")
	return nil
}
func startServer(ServerAddr string, r *gin.Engine) error {
	server := &http.Server{
		Addr:    ServerAddr,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Ошибка при запуске сервера: %v", err) // Меняем на log.Printf
		}
	}()

	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Ошибка при остановке сервера: %v\n", err)
		return err // Возвращаем ошибку
	}

	return nil
}
