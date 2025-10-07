package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/handler"
	"github.com/iamamatkazin/metrics.git/pkg/config"
)

func main() {
	cfg := config.New()
	router := handler.New().Router

	chSsignal := make(chan os.Signal, 1)
	signal.Notify(chSsignal, os.Interrupt, syscall.SIGTERM)

	tmRun := time.NewTimer(0)

loop:
	for {
		select {
		case <-chSsignal:
			slog.Info("Начало остановки сервера...")
			break loop

		case <-tmRun.C:
			go func() {
				http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), router)
			}()
		}
	}

	slog.Info("Выключение сервера")
}
