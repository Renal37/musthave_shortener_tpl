package main

import (
	"testing"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"honnef.co/go/tools/staticcheck"

	"github.com/Renal37/musthave_shortener_tpl.git/cmd/staticlint/analyzers"
)

// Тест на корректность создания мультичекера и работы анализаторов
func TestMultiCheckerFunctionality(t *testing.T) {
	// Добавляем анализаторы из staticcheck класса SA
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

	// Проверяем, что мультичекер успешно создается
	mc := multiChecker(myChecks...)
	if mc == nil {
		t.Fatal("multiChecker should not be nil")
	}

	// Проверяем, что у мультичекера правильное имя
	if mc.Name != "multiChecker" {
		t.Errorf("expected multiChecker name to be 'multiChecker', got %s", mc.Name)
	}
}
