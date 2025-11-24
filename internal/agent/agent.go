package agent

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/agent"
	pkghttp "github.com/iamamatkazin/metrics.git/pkg/http"
)

type Agent struct {
	client  *pkghttp.Client
	cfg     *agent.Config
	metrics map[string]map[string]float64
}

func New(cfg *agent.Config) *Agent {
	slog.Info("Запуск агента")
	return &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg),
		metrics: createMetrics(),
	}
}

func (a *Agent) Run(ctx context.Context) {
	pollCount := 0

	pollTicker := time.NewTicker(time.Second * time.Duration(a.cfg.PollInterval))
	defer pollTicker.Stop()

	reportTicker := time.NewTicker(time.Second * time.Duration(a.cfg.ReportInterval))
	defer reportTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-pollTicker.C:
			pollCount++
			a.poolMetrics(pollCount)

		case <-reportTicker.C:
			if err := a.sendMetricsOld(ctx); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}

			if err := a.sendMetricsBatch(ctx); err != nil {
				slog.Error("Ошибка отправки метрик на сервер:", slog.Any("error", err))
			}
		}
	}
}

func (a *Agent) sendMetricsOld(ctx context.Context) (err error) {
	urlBase := fmt.Sprintf("http://%s/update/", a.cfg.Address)

	for key, metrics := range a.metrics {
		for name, value := range metrics {
			url := fmt.Sprintf("%s%s/%s/%v", urlBase, key, name, value)

			if err := a.client.Post(ctx, url, "text/plain; charset=UTF-8", nil); err != nil {
				return err
			}

			// if err = a.sendSingleMetric(ctx, urlBase, key, name, value); err != nil {
			// 	return err
			// }
		}
	}

	return nil
}

// func (a *Agent) sendMetricNoBody(ctx context.Context, urlBase, key, name string, value float64) error {

// 	return nil
// }

// func (a *Agent) sendSingleMetric(ctx context.Context, urlBase, key, name string, value float64) error {
// 	metric := model.Metric{
// 		ID:    name,
// 		MType: key,
// 		Value: &value,
// 	}

// 	if err := a.client.Post(ctx, urlBase, "application/json", metric); err != nil {
// 		return err
// 	}

// 	return nil
// }

func (a *Agent) sendMetricsBatch(ctx context.Context) (err error) {
	urlBase := fmt.Sprintf("http://%s/updates/", a.cfg.Address)

	list := make([]model.Metric, 0, len(a.metrics[model.Gauge])+len(a.metrics[model.Counter]))
	for key, metrics := range a.metrics {
		for name, value := range metrics {
			list = append(list, model.Metric{ID: name, MType: key, Value: &value})
		}
	}

	if err := a.client.Post(ctx, urlBase, "application/json", list); err != nil {
		return err
	}

	return nil
}
