	package main

	import (
		"log"
		"net/http"
		_ "net/http/pprof"

		"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
		"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
		"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	)

	const (
		addr    = ":8080"  // адрес сервера
		maxSize = 10000000 // будем растить слайс до 10 миллионов элементов
	)

	func foo() {
		// полезная нагрузка
		for {
			var s []int
			for i := 0; i < maxSize; i++ {
				s = append(s, i)
			}
		}
	}
	func main() {
		addrConfig := config.InitConfig()
		storageInstance := storage.NewStorage()
		appInstance := app.NewApp(storageInstance, addrConfig)
		go foo()
		// Запуск pprof сервера на порту 6060
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()

		appInstance.Start()
		appInstance.Stop()
	}
