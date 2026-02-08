// Package custom реализует пользовательский анализатор.
package custom

import (
	"fmt"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Exit - кастомный анализатор запрета использования прямого вызова os.Exit.
var Exit = &analysis.Analyzer{
	Name: "errExit",
	Doc:  "запрет использования прямого вызова os.Exit",
	Run:  run,
}

// run - произоводит синтаксический анализ содержимого файла main.go пакета main на предмет
// запрета использования прямого вызова os.Exit.
func run(pass *analysis.Pass) (any, error) {
	fmt.Println(pass.Pkg.Name(), pass.Files)
	if pass.Pkg.Name() != "main" {
		return nil, nil
	}

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

	if mainFunc == nil || mainFunc.Body == nil {
		return nil, nil
	}

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

	return nil, nil
}
