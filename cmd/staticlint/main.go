package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/singlechecker"
	"honnef.co/go/tools/staticcheck"

	"github.com/Renal37/musthave_shortener_tpl.git/cmd/staticlint/analyzers"
)

func main() {
	// Добавляем анализаторы из staticcheck классов SA и других классов
	var staticcheckAnalyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		// Добавляем анализаторы из класса SA
		if v.Analyzer.Name[:2] == "SA" {
			staticcheckAnalyzers = append(staticcheckAnalyzers, v.Analyzer)
		}
		// Добавляем хотя бы один анализатор из других классов
		if v.Analyzer.Name[:1] == "S" && v.Analyzer.Name[:2] != "SA" {
			staticcheckAnalyzers = append(staticcheckAnalyzers, v.Analyzer)
		}
	}

	// Создаем список всех анализаторов, включая кастомный и публичные
	myChecks := []*analysis.Analyzer{
		inspect.Analyzer,           // Стандартный публичный анализатор
		shadow.Analyzer,            // Стандартный публичный анализатор
		nilness.Analyzer,           // Публичный анализатор на ваш выбор
		unusedresult.Analyzer,      // Публичный анализатор на ваш выбор
		analyzers.NoOsExitAnalyzer, // Кастомный анализатор
	}

	// Добавляем анализаторы из staticcheck
	myChecks = append(myChecks, staticcheckAnalyzers...)

	// Запускаем мультичекер
	singlechecker.Main(multiChecker(myChecks...))
}

// multiChecker объединяет несколько анализаторов для одновременного запуска
func multiChecker(analyzers ...*analysis.Analyzer) *analysis.Analyzer {
	return &analysis.Analyzer{
		Name: "multiChecker",
		Doc:  "выполняет несколько анализаторов",
		Run: func(pass *analysis.Pass) (interface{}, error) {
			for _, analyzer := range analyzers {
				_, err := analyzer.Run(pass)
				if err != nil {
					return nil, err
				}
			}
			return nil, nil
		},
	}
}
