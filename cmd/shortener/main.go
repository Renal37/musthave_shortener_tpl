package main

import (
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func main() {
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)
	appInstance.Start()
	appInstance.Stop()
}
