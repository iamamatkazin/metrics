package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/agent"
	"github.com/iamamatkazin/metrics.git/pkg/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.New()

	chSsignal := make(chan os.Signal, 1)
	signal.Notify(chSsignal, os.Interrupt, syscall.SIGTERM)

	chExit := make(chan struct{})
	tmRun := time.NewTimer(0)

loop:
	for {
		select {
		case <-chSsignal:
			slog.Info("Начало остановки агента...")
			cancel()

		case <-chExit:
			cancel()
			break loop

		case <-tmRun.C:
			go func() {
				agent.New(cfg).Run(ctx)
				close(chExit)
			}()
		}
	}

	slog.Info("Выключение агента")
}
