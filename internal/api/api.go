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
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RestAPI представляет собой структуру для REST API.
type RestAPI struct {
	Shortener *services.ShortenerService // Сервис для сокращения URL.
}

// StartRestAPI запускает REST API сервер.
// Он настраивает необходимые маршруты и middleware, и начинает прослушивание входящих запросов.
// Сервер завершает работу, когда переданный контекст отменяется.
//
// Параметры:
// - ctx: Контекст, используемый для управления жизненным циклом сервера.
// - ServerAddr: Адрес, на котором сервер будет прослушивать запросы.
// - BaseURL: Базовый URL для сервиса сокращения ссылок.
// - LogLevel: Уровень логирования для сервера.
// - db: Подключение к базе данных, используемое сервисом сокращения ссылок.
// - dbDNSTurn: Флаг, указывающий, включен ли DNS для базы данных.
// - storage: Хранилище, используемое сервисом сокращения ссылок.
//
// Возвращает ошибку, если сервер не удалось запустить или корректно завершить.
func StartRestAPI(ctx context.Context, ServerAddr, BaseURL, LogLevel string, db *repository.StoreDB, dbDNSTurn bool, storage *storage.Storage, EnableHTTPS bool, CertFile, KeyFile string) error {
	if err := logger.Initialize(LogLevel); err != nil {
		return fmt.Errorf("ошибка инициализации логгера: %w", err)
	}

	logger.Log.Info("Запуск сервера", zap.String("address", ServerAddr))
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

	// Создаем HTTP или HTTPS сервер
	srv := &http.Server{
		Addr:    ServerAddr,
		Handler: r,
	}

	go func() {
		var err error
		if EnableHTTPS {
			logger.Log.Info("Запуск HTTPS сервера")
			err = srv.ListenAndServeTLS(CertFile, KeyFile)
		} else {
			logger.Log.Info("Запуск HTTP сервера")
			err = srv.ListenAndServeTLS(CertFile, KeyFile)
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка при запуске сервера: %v\n", err)
		}
	}()

	// Ожидание завершения работы сервера
	<-ctx.Done()

	logger.Log.Info("Остановка сервера...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("ошибка при остановке сервера: %w", err)
	}

	logger.Log.Info("Сервер успешно остановлен")
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
