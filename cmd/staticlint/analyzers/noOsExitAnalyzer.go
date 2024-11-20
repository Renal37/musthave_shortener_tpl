package analyzers

import (
	"go/ast"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

// NoOsExitAnalyzer — кастомный анализатор, запрещающий использование os.Exit в main функции
var NoOsExitAnalyzer = &analysis.Analyzer{
	Name: "noOsExit",
	Doc:  "проверяет, что os.Exit не вызывается напрямую в main функции",
	Run:  run,
	Requires: []*analysis.Analyzer{
		inspect.Analyzer,
	},
}

// run выполняет проверку вызовов os.Exit
func run(pass *analysis.Pass) (interface{}, error) {
	// Проверяем, что анализируемый пакет — это main
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

	// Инициализируем инспектор для обхода AST
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	nodeFilter := []ast.Node{(*ast.CallExpr)(nil)}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		callExpr, _ := n.(*ast.CallExpr)
		if fun, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
			// Приведение типа fun.X к *ast.Ident для проверки имени пакета
			if pkgIdent, ok := fun.X.(*ast.Ident); ok && fun.Sel.Name == "Exit" && pkgIdent.Name == "os" {
				pass.Reportf(callExpr.Pos(), "не допускается прямой вызов os.Exit в функции main")
			}
		}
	})
	return nil, nil
}
