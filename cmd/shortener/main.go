package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func main() {
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	// Запуск pprof сервера, если флаг включен
	if addrConfig.EnablePprof {
		go func() {
			log.Println("Starting pprof server on :6060")
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	appInstance.Start()
	appInstance.Stop()
}
