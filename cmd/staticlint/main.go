package main

import (
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/singlechecker"
	"honnef.co/go/tools/staticcheck"

	"github.com/Renal37/musthave_shortener_tpl.git/cmd/staticlint/analyzers" 
)

func main() {
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
		analyzers.NoOsExitAnalyzer, // добавляем кастомный анализатор из нового пакета
	}
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
