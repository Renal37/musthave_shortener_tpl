package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
	"github.com/stretchr/testify/assert"
)

func TestAppInitialization(t *testing.T) {
	// Инициализируем конфигурацию и хранилище
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()
	appInstance := app.NewApp(storageInstance, addrConfig)

	assert.NotNil(t, appInstance)
}

func TestPprofServer(t *testing.T) {
	go func() {
		err := http.ListenAndServe("localhost:6060", nil)
		assert.NoError(t, err)
	}()
	time.Sleep(1 * time.Second)

	resp, err := http.Get("http://localhost:6060/debug/pprof/")
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
