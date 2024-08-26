package api

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/logger"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/middleware"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/store"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

type RestAPI struct {
	StructService *services.ShortenerService
}

func StartRestAPI(ServerAddr, BaseURL string, LogLevel string, db *store.StoreDB, dbDNSTurn bool, storage *storage.Storage) error {
	if err := logger.Initialize(LogLevel); err != nil {
		return err
	}
	logger.Log.Info("Running server", zap.String("address", ServerAddr))
	storageShortener := services.NewShortenerService(BaseURL, storage, db, dbDNSTurn)

	api := &RestAPI{
		StructService: storageShortener,
	}

	gin.SetMode(gin.ReleaseMode)
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