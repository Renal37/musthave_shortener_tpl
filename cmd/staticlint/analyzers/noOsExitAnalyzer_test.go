package analyzers_test

import (
	"testing"

	"github.com/Renal37/musthave_shortener_tpl.git/cmd/staticlint/analyzers"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestNoOsExitAnalyzer(t *testing.T) {
	// Создаем директорию testdata с Go-файлами для тестирования анализатора
	testdata := analysistest.TestData()

	// Запускаем анализатор на тестовых данных
	analysistest.Run(t, testdata, analyzers.NoOsExitAnalyzer, "a")
}
