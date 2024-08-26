package app

import (
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

func NewApp(storageInstance *storage.Storage, config *config.Config) *App {
	return &App{
		storageInstance: storageInstance,
		config:          config,
	}
}

func (a *App) Start() {
	db, err := store.InitDatabase(a.config.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	dbDNSTurn := true
	if a.config.DBPath == "" {
		err = dump.FillFromStorage(a.storageInstance, a.config.FilePath)
		if err != nil {
			log.Fatal(err)
		}
		dbDNSTurn = false
	}

	err = api.StartRestAPI(a.config.ServerAddr, a.config.BaseURL, a.config.LogLevel, db, dbDNSTurn, a.storageInstance)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) Stop() {
	if a.config.DBPath == "" {
		err := dump.Set(a.storageInstance, a.config.FilePath, a.config.BaseURL)
		if err != nil {
			log.Fatal(err)
		}
	}
}