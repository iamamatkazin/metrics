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

type request struct {
	url         string
	contentType string
	metric      any
}

type Agent struct {
	client  pkghttp.Clienter
	cfg     *agent.Config
	metrics map[string]map[string]float64
	jobs    chan request
}

func New(cfg *agent.Config) *Agent {
	slog.Info("Запуск агента")
	a := &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg),
		metrics: createMetrics(),
		jobs:    make(chan request, 100),
	}

	a.poolMetrics(1)
	a.poolGopsUtil()

	return a
}

func (a *Agent) Run(ctx context.Context) {
	for range a.cfg.RateLimit {
		go a.worker(ctx, a.jobs)
	}

	pollCount := 1

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
			go func(count int) {
				a.poolMetrics(count)
			}(pollCount)

			go a.poolGopsUtil()

		case <-reportTicker.C:
			a.sendMetricsOld()
			a.sendMetricsBatch()
		}
	}
}

func (a *Agent) sendMetricsOld() {
	urlBase := fmt.Sprintf("http://%s/update/", a.cfg.Address)

	for key, metrics := range a.metrics {
		for name, value := range metrics {
			url := fmt.Sprintf("%s%s/%s/%v", urlBase, key, name, value)

			select {
			case a.jobs <- request{url: url, contentType: "text/plain; charset=UTF-8", metric: nil}:
			default:
			}
		}
	}
}

func (a *Agent) sendMetricsBatch() {
	urlBase := fmt.Sprintf("http://%s/updates/", a.cfg.Address)

	list := make([]model.Metric, 0, len(a.metrics[model.Gauge])+len(a.metrics[model.Counter]))
	for key, metrics := range a.metrics {
		for name, value := range metrics {
			list = append(list, model.Metric{ID: name, MType: key, Value: &value})
		}
	}

	select {
	case a.jobs <- request{url: urlBase, contentType: "application/json", metric: list}:
	default:
	}
}
