package main

import (
	"github.com/Renal37/musthave_shortener_tpl.git/cmd/staticlint/analyzers"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/staticcheck"
	"testing"
)

// mockAnalyzer создает мок-анализатор для тестирования
func mockAnalyzer(name string, shouldFail bool) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: name,
		Run: func(pass *analysis.Pass) (interface{}, error) {
			if shouldFail {
				return nil, assert.AnError
			}
			return nil, nil
		},
	}
}

// Тест для multiChecker, проверяющий корректное выполнение анализаторов
func TestMultiChecker(t *testing.T) {
	// Создаем мок-анализаторы для тестирования
	analyzer1 := mockAnalyzer("mockAnalyzer1", false)
	analyzer2 := mockAnalyzer("mockAnalyzer2", false)

	// Создаем мультичекер с мок-анализаторами
	checker := multiChecker(analyzer1, analyzer2)

	// Проверяем, что мультичекер выполнит анализаторы без ошибок
	pass := &analysis.Pass{} // симулируем пустой pass, так как в данном тесте его неважно реализовывать
	_, err := checker.Run(pass)

	assert.NoError(t, err, "multiChecker должен выполняться без ошибок для успешных анализаторов")
}

// Тест для мультичекера с анализаторами из кода
func TestMainAnalyzers(t *testing.T) {
	// Собираем анализаторы из staticcheck класса SA
	var staticcheckAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		if v.Analyzer.Name[:2] == "SA" {
			staticcheckAnalyzers = append(staticcheckAnalyzers, v.Analyzer)
		}
	}

	// Создаем список всех анализаторов, включая кастомный
	myChecks := []*analysis.Analyzer{
		inspect.Analyzer,
		shadow.Analyzer,
		analyzers.NoOsExitAnalyzer,
	}
	myChecks = append(myChecks, staticcheckAnalyzers...)

	// Создаем мультичекер
	checker := multiChecker(myChecks...)

	// Проверяем, что мультичекер выполнит анализаторы без ошибок
	pass := &analysis.Pass{} // используем пустой pass для теста
	_, err := checker.Run(pass)

	assert.NoError(t, err, "multiChecker должен выполняться без ошибок для всех встроенных и кастомных анализаторов")
}
