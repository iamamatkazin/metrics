// Package main реализует агент сбора и отправки системных метрик.
// Агент представляет собой фоновую службу, которая периодически опрашивает
// системные метрики (CPU, память, GC статистика) и отправляет их на сервер
// сбора метрик через HTTP. Агент автоматически запускает два воркера:
//  1. Воркер сбора метрик - опрашивает систему с заданным интервалом
//  2. Воркер отправки метрик - отправляет собранные данные на сервер
package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iamamatkazin/metrics.git/internal/agent"
	"github.com/iamamatkazin/metrics.git/internal/common"
	aconfig "github.com/iamamatkazin/metrics.git/pkg/config/agent"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

// main - точка входа в приложение агента.
// Функция инициализирует конфигурацию, создает экземпляр агента
// и запускает фоновые процессы сбора и отправки метрик.
func main() {
	common.PrintBuild(buildVersion, buildDate, buildCommit)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := aconfig.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		return
	}

	serverPprof := &http.Server{
		Addr:    ":7171",
		Handler: nil,
	}

	go func() {
		slog.Info("Запуск сервера профилирования")
		if err := serverPprof.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера профилирования:", slog.Any("error", err))
		}
	}()

	go func() {
		if err := serverPprof.Shutdown(ctx); err != nil {
			slog.Error("Ошибка остановки сервера профилирования:", slog.Any("error", err))
		}

		<-ctx.Done()
		slog.Info("Начало остановки агента...")
	}()

	a, err := agent.New(cfg)
	if err != nil {
		slog.Error("Ошибка запуска агента:", slog.Any("error", err))
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		a.Start(ctx)
		a.Shutdown()
	}()

	go func() {
		defer wg.Done()

		for range cfg.RateLimit {
			a.Worker()
		}
	}()

	wg.Wait()

	slog.Info("Выключение агента")
}
