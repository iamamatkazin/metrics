package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/iamamatkazin/metrics.git/internal/agent"
	aconfig "github.com/iamamatkazin/metrics.git/pkg/config/agent"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := aconfig.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		os.Exit(2)
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

		a.Run(ctx)
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
