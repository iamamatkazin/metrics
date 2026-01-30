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
	aconfig "github.com/iamamatkazin/metrics.git/pkg/config/agent"
)

// main - точка входа в приложение агента.
// Функция инициализирует конфигурацию, создает экземпляр агента
// и запускает фоновые процессы сбора и отправки метрик.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := aconfig.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		return
	}

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit
		slog.Info("Начало остановки агента...")
		cancel()
	}()

	a := agent.New(cfg)
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

	go http.ListenAndServe(":7171", nil)

	wg.Wait()

	slog.Info("Выключение агента")
}
