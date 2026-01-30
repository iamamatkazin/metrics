// Package main реализует сервер сбора, хранения и выдачи метрик.
// Сервер представляет собой HTTP-службу, которая принимает метрики от агентов,
// сохраняет их в хранилище (in-memory, PostgreSQL или файловое) и предоставляет
// API для получения метрик. Поддерживает несколько типов хранилищ и механизмы
// аудита изменений.
package main

import (
	"context"
	"log/slog"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamamatkazin/metrics.git/internal/handler"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// main - точка входа в приложение сервера метрик.
// Функция загружает конфигурацию, инициализирует обработчики,
// запускает HTTP сервер и ожидает сигнала остановки.
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := sconfig.New()
	if err != nil {
		slog.Error("Ошибка чтения конфигурации:", slog.Any("error", err))
		return
	}

	app, err := handler.New(ctx, cfg)
	if err != nil {
		slog.Error("Ошибка создания сервера:", slog.Any("error", err))
		return
	}

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: app.Router,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	exit := make(chan struct{})

	go http.ListenAndServe(":7070", nil)

	go func() {
		slog.Info("Запуск сервера")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Ошибка запуска сервера:", slog.Any("error", err))
			close(exit)
		}
	}()

	select {
	case <-quit:
		app.Shutdown()

		if err := server.Shutdown(ctx); err != nil {
			slog.Error("Ошибка остановки сервера:", slog.Any("error", err))
		}
		cancel()

	case <-exit:
		return
	}

	slog.Info("Выключение сервера")
}
