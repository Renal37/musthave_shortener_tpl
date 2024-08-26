package app

import (
	"fmt"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/api"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/dump"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/Renal37/musthave_shortener_tpl.git/store"
)

type App struct {
	storageInstance *storage.Storage
	config          *config.Config
}

// NewApp создает новый экземпляр приложения с заданным хранилищем и конфигурацией
func NewApp(storageInstance *storage.Storage, config *config.Config) *App {
	return &App{
		storageInstance: storageInstance,
		config:          config,
	}
}

// Start запускает приложение: загружает данные из файла в хранилище или инициализирует базу данных
func (a *App) Start() {
	var err error
	var db *store.StoreDB

	// Проверка наличия DATABASE_DSN
	if a.config.DATABASE_DSN != "" {
		// Инициализируем базу данных
		db, err = store.InitDatabase(a.config.DATABASE_DSN)
		if err != nil {
			// Выводим ошибку, если не удалось инициализировать базу данных
			fmt.Printf("Ошибка при инициализации базы данных: %v\n", err)
			return
		}
	} else if a.config.FilePath != "" {
		// Если DATABASE_DSN не указан, используем файл
		err = dump.FillFromStorage(a.storageInstance, a.config.FilePath)
		if err != nil {
			fmt.Printf("Ошибка при заполнении хранилища из файла: %v\n", err)
			return
		}
	} else {
		// Если ни DATABASE_DSN, ни FilePath не указаны, используем память
		fmt.Println("Используется хранение данных в памяти.")
	}

	// Запускаем REST API
	err = api.StartRestAPI(a.config.ServerAddr, a.config.BaseURL, a.config.LogLevel, db, db != nil, a.storageInstance)
	if err != nil {
		// Выводим ошибку, если не удалось запустить API
		fmt.Printf("Ошибка при запуске REST API: %v\n", err)
	}
}

// Stop останавливает приложение: сохраняет данные из хранилища в файл, если используется файл
func (a *App) Stop() {
	if a.config.DATABASE_DSN == "" && a.config.FilePath != "" {
		err := dump.Set(a.storageInstance, a.config.FilePath, a.config.BaseURL)
		if err != nil {
			// Выводим ошибку, если не удалось сохранить данные
			fmt.Printf("Ошибка при сохранении данных в файл: %v\n", err)
		}
	}
}
