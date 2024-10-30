package logger

import (
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

// Initialize инициализирует глобальный логгер с заданным уровнем логирования
func Initialize(level string) error {
	// Парсим уровень логирования из строки
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err // Возвращаем ошибку, если уровень логирования некорректный
	}

	// Создаем конфигурацию для логгера в режиме "production"
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl // Устанавливаем уровень логирования

	// Строим новый логгер с указанной конфигурацией
	zl, err := cfg.Build()
	if err != nil {
		return err // Возвращаем ошибку, если не удалось создать логгер
	}

	// Устанавливаем глобальный логгер
	Log = zl.Sugar()
	return nil
}
