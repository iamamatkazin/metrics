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
	cfg     *config.Config
	metrics map[string]map[string]any
}

func New(cfg *config.Config) *Agent {
	slog.Info("Запуск агента")
	return &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg),
		metrics: createMetrics(),
	}
}

func (a *Agent) Run(ctx context.Context) {
	pollCount := 0

	tiPool := time.NewTicker(a.cfg.Client.PollInterval)
	defer tiPool.Stop()

	tiReport := time.NewTicker(a.cfg.Client.ReportInterval)
	defer tiReport.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-tiPool.C:
			pollCount++
			a.poolMetrics(pollCount)

		case <-tiReport.C:
			a.reportMetrics(ctx)
		}

	}
}

func (a *Agent) reportMetrics(ctx context.Context) (err error) {
	urlBase := fmt.Sprintf("http://%s:%d/update", a.cfg.Server.Host, a.cfg.Server.Port)

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
