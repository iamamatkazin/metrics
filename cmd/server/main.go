package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamamatkazin/metrics.git/internal/handler"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := server.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: handler.New().Router,
	}
	exit := make(chan struct{})

	go func() {
		slog.Info("Запуск сервера")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера:", slog.Any("error", err))
			close(exit)
		}
	}()

	select {
	case <-quit:
		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Ошибка остановки сервера:", slog.Any("error", err))
		}
		cancel()

	case <-exit:
		os.Exit(2)
	}

	slog.Info("Выключение сервера")
}
