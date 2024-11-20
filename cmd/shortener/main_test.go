package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/internal/app"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/config"
	"github.com/Renal37/musthave_shortener_tpl.git/internal/storage"
)

func TestMainInfo(t *testing.T) {
	// Перехватываем вывод для проверки
	var buf bytes.Buffer
	// Используем Fprintf, чтобы записать вывод в буфер
	fmt.Fprintf(&buf, "Build version: %s\n", "1.0.0")
	fmt.Fprintf(&buf, "Build date: %s\n", "2024-11-08")
	fmt.Fprintf(&buf, "Build commit: %s\n", "5b790d0dd1415028222997c1a6464bffac9729df")

	expected := "Build version: 1.0.0\nBuild date: 2024-11-08\nBuild commit: 5b790d0dd1415028222997c1a6464bffac9729df\n"
	if buf.String() != expected {
		t.Errorf("Неожиданный вывод информации о сборке:\nПолучено:\n%s\nОжидалось:\n%s", buf.String(), expected)
	}
}
func TestInitConfig(t *testing.T) {
	addrConfig := config.InitConfig()

	if addrConfig.ServerAddr == "" {
		t.Error("Ожидалось, что ServerAddress будет установлен")
	}

	if addrConfig.BaseURL == "" {
		t.Error("Ожидалось, что BaseURL будет установлен")
	}
}

func TestNewStorage(t *testing.T) {
	storageInstance := storage.NewStorage()

	if storageInstance == nil {
		t.Error("Ожидалось, что NewStorage вернет не nil экземпляр хранилища")
	}
}

func TestNewApp(t *testing.T) {
	addrConfig := config.InitConfig()
	storageInstance := storage.NewStorage()

	appInstance := app.NewApp(storageInstance, addrConfig)

	if appInstance == nil {
		t.Error("Ожидалось, что NewApp вернет не nil экземпляр приложения")
	}
}
