package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
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
	sync.RWMutex
}

func New(cfg *agent.Config) *Agent {
	slog.Info("Запуск агента")
	a := &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg.Timeout),
		metrics: createMetrics(),
		jobs:    make(chan request, 100),
	}

	return a
}

func (a *Agent) Run(ctx context.Context) {
	pollCount := 1

	a.poolMetrics(pollCount)
	a.poolGopsUtil()

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
	a.RLock()
	defer a.RUnlock()

	urlBase := fmt.Sprintf("http://%s/update/", a.cfg.Address)

	for key, metrics := range a.metrics {
		for name, value := range metrics {
			url := fmt.Sprintf("%s%s/%s/%v", urlBase, key, name, value)

			select {
			case a.jobs <- request{url: url, contentType: "text/plain; charset=UTF-8", metric: nil}:
			default:
				slog.Info("Нет свободного канала для обработки метрики:", slog.Any(name, value))
			}
		}
	}
}

func (a *Agent) sendMetricsBatch() {
	a.RLock()
	defer a.RUnlock()

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
		slog.Info("Нет свободного канала для обработки списка метрик.")
	}
}

func (a *Agent) Shutdown() {
	close(a.jobs)
}
