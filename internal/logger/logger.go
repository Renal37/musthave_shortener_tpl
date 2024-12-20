package logger

import (
	"go.uber.org/zap"
)

// Log представляет собой глобальный логгер, который используется в приложении.
var Log *zap.SugaredLogger

// Initialize инициализирует глобальный логгер с заданным уровнем логирования.
//
// level - строка, представляющая уровень логирования (например, "info", "debug").
// Возвращает ошибку, если уровень логирования некорректен или если не удалось создать логгер.
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
