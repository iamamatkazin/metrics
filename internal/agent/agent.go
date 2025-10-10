package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/iamamatkazin/metrics.git/pkg/config"
	pkghttp "github.com/iamamatkazin/metrics.git/pkg/http"
)

type Agent struct {
	client  *pkghttp.Client
	cfg     *config.Client
	metrics map[string]map[string]any
}

func New(cfg *config.Client) *Agent {
	slog.Info("Запуск агента")
	return &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg),
		metrics: createMetrics(),
	}
}

func (a *Agent) Run(ctx context.Context) {
	pollCount := 0

	pollTicker := time.NewTicker(a.cfg.PollInterval)
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(a.cfg.ReportInterval)
	defer reportTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-pollTicker.C:
			pollCount++
			a.poolMetrics(pollCount)

		case <-reportTicker.C:
			a.reportMetrics(ctx)
		}
	}
}

func (a *Agent) reportMetrics(ctx context.Context) (err error) {
	urlBase := fmt.Sprintf("http://%s/update", a.cfg.ServerAddress)

	for key, metrics := range a.metrics {
		for name, value := range metrics {
			url := fmt.Sprintf("%s/%s/%s/%v", urlBase, key, name, value)

			if err = a.client.Post(ctx, url, "text/plain; charset=UTF-8"); err != nil {
				return err
			}
		}
	}

	return nil
}
