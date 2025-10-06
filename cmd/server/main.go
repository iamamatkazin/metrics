package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/iamamatkazin/metrics.git/pkg/config"
)

func updateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Сервер поддерживает только POST запросы", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/update/", updateMetrics)

	cfg := config.New()

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
				http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
			}()
		}
	}

	slog.Info("Выключение сервера")
}
