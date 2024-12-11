package app

import (
	"context"
	"fmt"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api/grpc"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api/rest"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/services"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
	"golang.org/x/sync/errgroup"
)

// App represents the application structure containing storage and configuration.
type App struct {
	storageInstance  *storage.Storage           // Pointer to storage
	servicesInstance *services.ShortenerService // Pointer to services
	config           *config.Config             // Pointer to configuration
	fillFromStorage  func(*storage.Storage, string) error
	set              func(*storage.Storage, string) error
}

// NewApp creates a new instance of the application with the given storage and configuration.
func NewApp(storageInstance *storage.Storage, servicesInstance *services.ShortenerService, config *config.Config) *App {
	return &App{
		storageInstance:  storageInstance,
		servicesInstance: servicesInstance,
		config:           config,
		fillFromStorage:  dump.FillFromStorage,
		set:              dump.Set,
	}
}

// Start запускает приложение: загружает данные из файла в хранилище и запускает REST API.
func (a *App) Start(ctx context.Context) error {
	// Инициализируем базу данных
	db, err := repository.InitDatabase(a.config.DBPath)
	if err != nil {
		fmt.Printf("Ошибка при инициализации базы данных: %v\n", err)
		return err
	}

	dbDNSTurn := true
	if a.UseDatabase() {
		err = dump.FillFromStorage(a.storageInstance, a.config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка при заполнении хранилища: %v\n", err)
			return err
		}
		dbDNSTurn = false
	}

	var eg errgroup.Group
	// Запускаем REST API с контекстом
	eg.Go(func() error {
		err := rest.StartRestAPI(
			ctx,
			a.config.ServerAddr,
			a.config.BaseURL,
			a.config.LogLevel,
			db,
			dbDNSTurn,
			a.storageInstance,
			a.config.EnableHTTPS,
			a.config.CertFile,
			a.config.KeyFile,
		)
		return err
	})

	// Запускаем gRPC сервер
	eg.Go(func() error {
		err := grpc.StartGRPCServer(ctx, a.servicesInstance, ":50051")
		return err
	})

	eg.Go(func() error {
		// Ожидание завершения контекста
		<-ctx.Done()
		if err := a.Stop(); err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		fmt.Printf("Ошибка при запуске приложения: %v\n", err)
		return err
	}
	return nil
}

// UseDatabase возвращает true, если приложение использует базу данных.
func (a *App) UseDatabase() bool {
	return a.config.DBPath == ""
}

// Stop останавливает приложение: сохраняет данные из хранилища в файл.
func (a *App) Stop() error {
	fmt.Println("Сохраняем данные перед завершением работы...")
	if a.UseDatabase() {
		err := a.set(a.storageInstance, a.config.FilePath) // Используем a.set вместо dump.Set
		if err != nil {
			return err
		}
		fmt.Println("Данные успешно сохранены.")
	}
	return nil
}
