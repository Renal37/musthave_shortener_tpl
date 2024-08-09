package app

import (
	"fmt"
	"log"
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

// Start запускает приложение: загружает данные из файла в хранилище и запускает REST API
func (a *App) Start() {
	// Инициализируем базу данных
	db, err := store.InitDatabase(a.config.DBPath)
	if err != nil {
		// Выводим ошибку, если не удалось инициализировать базу данных
		fmt.Printf("Ошибка при инициализации базы данных: %v\n", err)
		return
	}

	dbDNSTurn := true
	if a.config.DBPath == "" {
		// Заполняем хранилище данными из файла
		err = dump.FillFromStorage(a.storageInstance, a.config.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		dbDNSTurn = false
	}

	// Запускаем REST API
	err = api.StartRestAPI(a.config.ServerAddr, a.config.BaseURL, a.config.LogLevel, db, dbDNSTurn, a.storageInstance)
	if err != nil {
		// Выводим ошибку, если не удалось запустить API
		fmt.Printf("Ошибка при запуске REST API: %v\n", err)
	}
}

// Stop останавливает приложение: сохраняет данные из хранилища в файл
func (a *App) Stop() {
	// Сохраняем данные из хранилища в файл
	if a.config.DBPath == "" {
		err := dump.Set(a.storageInstance, a.config.FilePath, a.config.BaseURL)
		if err != nil {
			// Выводим ошибку, если не удалось сохранить данные
			fmt.Printf("Ошибка при сохранении данных: %v\n", err)
		}
	}
}