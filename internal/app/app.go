package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/repository"
)

// App представляет собой структуру приложения, содержащую хранилище и конфигурацию.
type App struct {
	storageInstance *storage.Storage // Указатель на хранилище
	config          *config.Config   // Указатель на конфигурацию
	apiShutdown     func() error     // Функция для завершения REST API
}

// NewApp создает новый экземпляр приложения с заданным хранилищем и конфигурацией.
func NewApp(storageInstance *storage.Storage, config *config.Config) *App {
	return &App{
		storageInstance: storageInstance,
		config:          config,
	}
}

// Start запускает приложение: загружает данные из файла в хранилище и запускает REST API.
func (a *App) Start() error {
	// Инициализируем базу данных
	db, err := repository.InitDatabase(a.config.DBPath)
	if err != nil {
		fmt.Printf("Ошибка при инициализации базы данных: %v\n", err)
		return nil
	}

	dbDNSTurn := true
	if a.UseDatabase() {
		err = dump.FillFromStorage(a.storageInstance, a.config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка при заполнении хранилища: %v\n", err)
			return nil
		}
		dbDNSTurn = false
	}

	// Контекст для завершения работы приложения
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Канал для получения системных сигналов
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Канал для завершения API
	apiDone := make(chan error, 1)

	// Запускаем REST API
	go func() {
		err := api.StartRestAPI(a.config.ServerAddr, a.config.BaseURL, a.config.LogLevel, db, dbDNSTurn, a.storageInstance)
		apiDone <- err
	}()

	// Обработка сигналов завершения
	go func() {
		sig := <-signalChan
		fmt.Printf("Получен сигнал: %v. Завершаем работу...\n", sig)
		cancel()
		a.Stop()
	}()

	// Ожидание завершения API или получения ошибки
	select {
	case <-ctx.Done():
		fmt.Println("Контекст завершён")
	case err := <-apiDone:
		if err != nil {
			fmt.Printf("Ошибка при запуске REST API: %v\n", err)
		}
	}

	return nil
}

// UseDatabase возвращает true, если приложение использует базу данных.
func (a *App) UseDatabase() bool {
	return a.config.DBPath == ""
}

// Stop останавливает приложение: сохраняет данные из хранилища в файл.
func (a *App) Stop() {
	fmt.Println("Сохраняем данные перед завершением работы...")
	if a.UseDatabase() {
		err := dump.Set(a.storageInstance, a.config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка при сохранении данных: %v\n", err)
		} else {
			fmt.Println("Данные успешно сохранены.")
		}
	}
}
