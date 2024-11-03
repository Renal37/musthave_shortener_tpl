package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

var profilingEnabled = flag.Bool("profiling", false, "Enable profiling server")

func main() {
	flag.Parse()

	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	if *profilingEnabled {
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	appInstance.Start()
	appInstance.Stop()
}
