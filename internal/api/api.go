package api

// Импортируем необходимые пакеты
import (
	"context"
	"fmt"
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

// Определяем структуру RestAPI, которая будет хранить ссылку на ShortenerService
type RestAPI struct {
	StructService *services.ShortenerService
}

// Функция StartRestAPI запускает REST API сервер
func StartRestAPI(ServerAddr, BaseURL string, LogLevel string, db *store.StoreDB, dbDNSTurn bool, storage *storage.Storage) error {
	// Инициализируем логгер с указанным уровнем логирования
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}
	// Выводим сообщение об успешном запуске сервера с указанным адресом
	logger.Log.Info("Запуск сервера", zap.String("адрес", ServerAddr))
	// Создаем новый ShortenerService с указанными параметрами
	storageShortener := services.NewShortenerService(BaseURL, storage, db, dbDNSTurn)
	// Создаем новый экземпляр RestAPI с указанным ShortenerService
	api := &RestAPI{
		StructService: storageShortener,
	}
	// Устанавливаем режим работы Gin на ReleaseMode
	gin.SetMode(gin.ReleaseMode)
	// Создаем новый Gin router
	r := gin.Default()
	// Используем Gin recovery middleware
	r.Use(gin.Recovery())
	// Используем LoggerMiddleware с указанным логгером
	r.Use(middleware.LoggerMiddleware(logger.Log))
	// Используем CompressMiddleware
	r.Use(middleware.CompressMiddleware())
	r.Use(middleware.AuthorizationMiddleware())
	// Устанавливаем маршруты для API с помощью функции api.setRoutes
	api.setRoutes(r)
	// Создаем новый HTTP сервер с адресом ":8080" и указанным router в качестве Handler
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	// Создаем канал для получения OS сигналов
	quit := make(chan os.Signal, 1)
	// Регистрируем os.Interrupt сигнал в канал quit
	signal.Notify(quit, os.Interrupt)
	// Запускаем горутину для запуска сервера на указанном адресе
	go func() {
		err := r.Run(ServerAddr)
		if err != nil {
			fmt.Println("Не удалось запустить браузер")
		}
	}()
	// Ждем получения сигнала в канал quit
	<-quit
	// Создаем контекст с 5-секундным таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// Отменяем контекст при завершении функции
	defer cancel()
	// Грациозно остановляем сервер с помощью контекста
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("Ошибка при остановке сервера: %v\n", err)
	}
	// Если сервер успешно запущен и остановлен, возвращаем nil
	return nil
}
