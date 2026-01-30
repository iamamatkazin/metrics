// Package main реализует статический анализатор Go кода с набором проверок.
// Staticlint объединяет множество стандартных анализаторов Go (из пакета
// golang.org/x/tools/go/analysis/passes) и дополнительные анализаторы
// из пакета honnef.co/go/tools (staticcheck, simple, stylecheck, quickfix).
// + Кастомный анализатор:
//   - custom.Exit - проверяет отсутствие прямых вызовов os.Exit
package main

import (
	"github.com/iamamatkazin/metrics.git/cmd/staticlint/custom"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/asmdecl"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildssa"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/cgocall"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/ctrlflow"
	"golang.org/x/tools/go/analysis/passes/deepequalerrors"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/fieldalignment"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/framepointer"
	"golang.org/x/tools/go/analysis/passes/gofix"
	"golang.org/x/tools/go/analysis/passes/hostport"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/inline"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/reflectvaluecompare"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/slog"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"golang.org/x/tools/go/analysis/passes/waitgroup"
	"honnef.co/go/tools/analysis/lint"
	"honnef.co/go/tools/quickfix"
	"honnef.co/go/tools/simple"
	"honnef.co/go/tools/staticcheck"
	"honnef.co/go/tools/stylecheck"
)

// main - точка входа в приложение статического анализатора.
// Функция собирает все анализаторы и запускает multichecker
// для проверки Go кода.
func main() {
	passes := getAnalysesPasses()
	saAnalyzers := getAnalysesSA()
	qfAnalyzers := getAnalysesQuickfix()
	sAnalyzers := getAnalysesSimple()
	stAnalyzers := getAnalysesStylecheck()

	analyses := make([]*analysis.Analyzer, 0,
		1+len(passes)+len(saAnalyzers)+len(qfAnalyzers)+len(sAnalyzers)+len(stAnalyzers))

	analyses = append(analyses, custom.Exit)
	analyses = append(analyses, passes...)
	analyses = append(analyses, saAnalyzers...)
	analyses = append(analyses, qfAnalyzers...)
	analyses = append(analyses, sAnalyzers...)
	analyses = append(analyses, stAnalyzers...)

	multichecker.Main(
		analyses...,
	)
}

// getAnalysesPasses - формирование массива анализаторов из analysis.
func getAnalysesPasses() []*analysis.Analyzer {
	return []*analysis.Analyzer{
		appends.Analyzer,
		asmdecl.Analyzer,
		assign.Analyzer,
		atomic.Analyzer,
		atomicalign.Analyzer,
		bools.Analyzer,
		buildssa.Analyzer,
		buildtag.Analyzer,
		cgocall.Analyzer,
		composite.Analyzer,
		copylock.Analyzer,
		ctrlflow.Analyzer,
		deepequalerrors.Analyzer,
		defers.Analyzer,
		directive.Analyzer,
		errorsas.Analyzer,
		fieldalignment.Analyzer,
		findcall.Analyzer,
		framepointer.Analyzer,
		gofix.Analyzer,
		hostport.Analyzer,
		httpmux.Analyzer,
		httpresponse.Analyzer,
		ifaceassert.Analyzer,
		inline.Analyzer,
		inspect.Analyzer,
		loopclosure.Analyzer,
		lostcancel.Analyzer,
		nilfunc.Analyzer,
		nilness.Analyzer,
		pkgfact.Analyzer,
		printf.Analyzer,
		reflectvaluecompare.Analyzer,
		shadow.Analyzer,
		shift.Analyzer,
		sigchanyzer.Analyzer,
		slog.Analyzer,
		sortslice.Analyzer,
		stdmethods.Analyzer,
		stdversion.Analyzer,
		stringintconv.Analyzer,
		structtag.Analyzer,
		testinggoroutine.Analyzer,
		tests.Analyzer,
		timeformat.Analyzer,
		unmarshal.Analyzer,
		unreachable.Analyzer,
		unsafeptr.Analyzer,
		unusedresult.Analyzer,
		unusedwrite.Analyzer,
		usesgenerics.Analyzer,
		waitgroup.Analyzer,
	}
}

// getAnalysesSA - формирование массива анализаторов из staticcheck.
func getAnalysesSA() []*analysis.Analyzer {
	var analyzers []*analysis.Analyzer
	for _, v := range staticcheck.Analyzers {
		analyzers = append(analyzers, v.Analyzer)
	}

	return analyzers
}

// getAnalysesQuickfix - формирование массива анализаторов из quickfix.
func getAnalysesQuickfix() []*analysis.Analyzer {
	// создаем перечень для анализаторов из quickfix
	target := map[string]bool{
		"QF1001": true,
		"QF1002": true,
		"QF1003": true,
	}

	return fillAnalyzers(target, quickfix.Analyzers)
}

// getAnalysesSimple - формирование массива анализаторов из simple.
func getAnalysesSimple() []*analysis.Analyzer {
	// создаем перечень для анализаторов из simple
	target := map[string]bool{
		"S1001": true,
		"S1002": true,
		"S1003": true,
	}

	return fillAnalyzers(target, simple.Analyzers)
}

// getAnalysesStylecheck - формирование массива анализаторов из stylecheck.
func getAnalysesStylecheck() []*analysis.Analyzer {
	// создаем перечень для анализаторов из stylecheck
	target := map[string]bool{
		"ST1001": true,
		"ST1003": true,
		"ST1005": true,
	}

	return fillAnalyzers(target, stylecheck.Analyzers)
}

// fillAnalyzers - формирует пользовательский слайс анализаторов.
func fillAnalyzers(target map[string]bool, source []*lint.Analyzer) []*analysis.Analyzer {
	analyzers := make([]*analysis.Analyzer, 0, len(target))

	for _, v := range source {
		if target[v.Analyzer.Name] {
			analyzers = append(analyzers, v.Analyzer)
		}
	}

	return analyzers
}
