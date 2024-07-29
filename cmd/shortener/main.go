package main

import (
	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
)

func main() {
	addrConfig := config.InitConfig()
	app.Start(addrConfig)
}