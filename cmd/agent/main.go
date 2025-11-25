package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go func() {
		agent.New(cfg).Run(ctx)
		close(exit)
	}()

	select {
	case <-quit:
		slog.Info("Начало остановки агента...")
		cancel()
	case <-exit:
		cancel()
	}

	slog.Info("Выключение агента")
}
