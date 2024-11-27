package api

import (
	"context"
	"log"
	"net/http"
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
//   - ctx: Контекст для управления жизненным циклом сервера.
//   - ServerAddr: Адрес для запуска сервера, например, ":8080".
//   - BaseURL: Базовый URL сервиса сокращения ссылок.
//   - LogLevel: Уровень логирования, например, "info" или "debug".
//   - db: Объект StoreDB для хранения данных в базе.
//   - dbDNSTurn: Логический флаг, указывающий, использовать ли DNS-трансляцию для базы данных.
//   - storage: Объект Storage для хранения данных.
//
// Возвращает указатель на HTTP-сервер и функцию завершения API.
func StartRestAPI(ctx context.Context, ServerAddr, BaseURL string, LogLevel string, db *repository.StoreDB, dbDNSTurn bool, storage *storage.Storage) (*http.Server, func() error) {
	if err := logger.Initialize(LogLevel); err != nil {
		log.Fatalf("Ошибка инициализации логгера: %v", err)
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

	server := &http.Server{
		Addr:    ServerAddr,
		Handler: r,
	}

	// Горутинa для остановки сервера по завершению контекста
	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Log.Error("Ошибка остановки сервера", zap.Error(err))
		}
	}()

	return server, server.Close
}
