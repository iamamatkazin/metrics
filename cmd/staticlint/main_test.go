package main

import (
	"github.com/iamamatkazin/metrics.git/cmd/staticlint/custom"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
)

func Example() {
	// формируем слайс из всех анализаторов Passes
	passes := getAnalysesPasses()

	// формируем слайс из всех анализаторов SA
	saAnalyzers := getAnalysesSA()

	// формируем слайс из трех анализаторов (QF1001, QF1002, QF1003) Quickfix
	qfAnalyzers := getAnalysesQuickfix()

	// формируем слайс из трех анализаторов (S1001, S1002, S1003) Simple
	sAnalyzers := getAnalysesSimple()

	// формируем слайс из трех анализаторов (ST1001, ST1003, ST1005) Stylecheck
	stAnalyzers := getAnalysesStylecheck()

	// инициализируем общий слайс для всех анализаторов
	analyses := make([]*analysis.Analyzer, 0,
		1+len(passes)+len(saAnalyzers)+len(qfAnalyzers)+len(sAnalyzers)+len(stAnalyzers))

	// заполняем общий слайс всех анализаторов, анализаторами из каждой группы
	analyses = append(analyses, custom.Exit)
	analyses = append(analyses, passes...)
	analyses = append(analyses, saAnalyzers...)
	analyses = append(analyses, qfAnalyzers...)
	analyses = append(analyses, sAnalyzers...)
	analyses = append(analyses, stAnalyzers...)

	// запускаем multichecker
	multichecker.Main(
		analyses...,
	)
}
