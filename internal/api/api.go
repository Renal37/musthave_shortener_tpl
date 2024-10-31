package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
func StartRestAPI(ServerAddr, BaseURL string, LogLevel string, db *repository.StoreDB, dbDNSTurn bool, storage *storage.Storage) error {
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

	api.setRoutes(r) // Настройка маршрутов API

	server := &http.Server{
		Addr:    ServerAddr,
		Handler: r,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v", err)
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
