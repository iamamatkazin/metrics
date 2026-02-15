package custom

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

func Example() {
	var pass *analysis.Pass
	// Если передам не пакет main, то выходим
	if pass.Pkg.Name() != "main" {
		// выход из обработки
	}

	// в цикле перебираем все файлы пакета main и ищем в них функцию main
	var mainFunc *ast.FuncDecl
loop:
	for _, file := range pass.Files {
		for _, decl := range file.Decls {
			if funcDecl, ok := decl.(*ast.FuncDecl); ok && funcDecl.Recv == nil && funcDecl.Name.Name == "main" {
				mainFunc = funcDecl
				break loop
			}
		}
	}

	// если функция не найдена или не имеет тела, то завершаем обработку
	if mainFunc == nil || mainFunc.Body == nil {
		// выход из обработки
	}

	// обходим все узлы функции main
	ast.Inspect(mainFunc.Body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		var ident *ast.Ident
		switch fun := call.Fun.(type) {
		case *ast.Ident:
			ident = fun
		case *ast.SelectorExpr:
			ident = fun.Sel // Sel — это *ast.Ident для имени функции
		default:
			return true
		}

		// Получаем объект через TypesInfo.Uses
		obj, ok := pass.TypesInfo.Uses[ident]
		if !ok {
			return true
		}

		// Проверяем, что это именно функция os.Exit
		if fn, ok := obj.(*types.Func); ok {
			if pkg := fn.Pkg(); pkg != nil && pkg.Path() == "os" && fn.Name() == "Exit" {
				pass.Reportf(call.Pos(), "прямой вызов os.Exit() внутри функции main запрещён")
			}
		}

		return true
	})

	// завершение анализа
}
