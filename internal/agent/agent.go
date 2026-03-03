// Package agent предоставляет функционал для сбора и отправки метрик.
// Реализует фоновый воркер, который периодически опрашивает системные метрики,
// а затем отправляет их на сервер сбора метрик через HTTP.
package agent

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/grpc/client"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/agent"
	pkghttp "github.com/iamamatkazin/metrics.git/pkg/http"
)

// request представляет запрос метрики для внутренней передачи данных.
type request struct {
	metric      any
	url         string
	contentType string
}

// Agent - основная структура агента сбора метрик.
// Содержит конфигурацию, HTTP клиент и хранилище метрик.
type Agent struct {
	client  pkghttp.Clienter
	cfg     *agent.Config
	metrics map[string]map[string]float64
	jobs    chan request
	grpc    *client.Metrics
	sync.RWMutex
}

// New создает новый экземпляр агента с заданной конфигурацией.
// Агент начнет сбор метрик согласно настройкам конфигурации.
func New(cfg *agent.Config) (*Agent, error) {
	slog.Info("Запуск агента")

	grpc, err := client.New(cfg)
	if err != nil {
		return nil, err
	}

	a := &Agent{
		cfg:     cfg,
		client:  pkghttp.New(cfg.Timeout),
		metrics: createMetrics(),
		jobs:    make(chan request, 100),
		grpc:    grpc,
	}

	return a, nil
}

// Start запускает процесс сбора и отправки метрик.
// Запускает два цикла на основе тикеров: один для опроса метрик и один для отправки.
func (a *Agent) Start(ctx context.Context) {
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

// sendMetricsOld отправляет метрики через параметры URL (устаревший метод).
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

// sendMetricsBatch отправляет несколько метрик в одном JSON запросе.
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

	if err := a.sendMetricsForGrpc(list); err != nil {
		slog.Error("Ошибка отправка метрик через grpc:", slog.Any("error", err))
	}
}

func (a *Agent) sendMetricsForGrpc(list []model.Metric) error {
	ctx, cancel := context.WithTimeout(context.Background(), a.cfg.Timeout)
	defer cancel()

	if err := a.grpc.UpdateMetrics(ctx, list); err != nil {
		return err
	}

	return nil
}

// Shutdown коректно останавливает агента и закрывает канал jobs.
func (a *Agent) Shutdown() {
	close(a.jobs)
	a.grpc.Close()
}
