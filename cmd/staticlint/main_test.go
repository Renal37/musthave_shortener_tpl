package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
)

func TestMultiChecker(t *testing.T) {
	// Создаем фиктивный анализатор, который всегда успешно выполняется
	successAnalyzer := &analysis.Analyzer{
		Name: "successAnalyzer",
		Doc:  "Успешный анализатор для тестирования",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return nil, nil
		},
	}

	// Создаем фиктивный анализатор, который всегда возвращает ошибку
	failingAnalyzer := &analysis.Analyzer{
		Name: "failingAnalyzer",
		Doc:  "Анализатор с ошибкой для тестирования",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			return nil, assert.AnError
		},
	}

	// Тестируем multiChecker с успешным анализатором
	multi := multiChecker(successAnalyzer)
	_, err := multi.Run(nil)
	assert.NoError(t, err, "Ожидалось отсутствие ошибки от успешного анализатора")

	// Тестируем multiChecker с анализатором, который возвращает ошибку
	multi = multiChecker(failingAnalyzer)
	_, err = multi.Run(nil)
	assert.Error(t, err, "Ожидалась ошибка от анализатора с ошибкой")

	// Тестируем multiChecker с несколькими анализаторами
	multi = multiChecker(successAnalyzer, failingAnalyzer)
	_, err = multi.Run(nil)
	assert.Error(t, err, "Ожидалась ошибка, так как один из анализаторов возвращает ошибку")

	// Тестируем multiChecker с двумя успешными анализаторами
	multi = multiChecker(successAnalyzer, successAnalyzer)
	_, err = multi.Run(nil)
	assert.NoError(t, err, "Ожидалось отсутствие ошибки, так как все анализаторы успешны")
}
