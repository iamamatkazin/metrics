// Package memstorage предоставляет in-memory хранилище для метрик с опциональным
// периодическим файловым бэкапом. Реализует интерфейс Storager и используется
// как основной слой доступа к метрикам.
package memstorage

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/repository/filestorage"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Storage реализует интерфейс Storager для хранения метрик в памяти.
type Storage struct {
	metrics  map[string]*model.Metric
	cfg      *server.Config
	fileStor *filestorage.Storage
	sync.RWMutex
}

// New создает новое in-memory хранилище для метрик.
func New(ctx context.Context, cfg *server.Config) (*Storage, error) {
	fileStor, err := filestorage.New(cfg)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]*model.Metric)
	if cfg.Restore {
		if err = fileStor.LoadDump(metrics); err != nil {
			return nil, err
		}
	}

	s := &Storage{
		cfg:      cfg,
		fileStor: fileStor,
		metrics:  metrics,
	}

	go func() {
		s.saveDump(ctx)
	}()

	return s, nil
}

// GetMetric получает метрику по её идентификатору.
func (s *Storage) GetMetric(_ context.Context, id string) (*model.Metric, error) {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.metrics[id]
	if !ok {
		return nil, nil
	}

	return val, nil
}

// UpdateMetric сохраняет новую метрику или изменяет существующую.
func (s *Storage) UpdateMetric(_ context.Context, metric *model.Metric) error {
	s.Lock()

	val, ok := s.metrics[metric.ID]
	if !ok {
		s.metrics[metric.ID] = metric
	} else {
		if val.MType == model.Gauge {
			val.Value = metric.Value
		} else {
			if metric.Delta != nil {
				delta := *val.Delta + *metric.Delta
				val.Delta = &delta
			} else if metric.Value != nil {
				delta := *val.Delta + int(*metric.Value)
				val.Delta = &delta
			}
		}
	}
	s.Unlock()

	if s.cfg.StoreInterval == 0 {
		s.saveToFile()
	}

	return nil
}

// ListMetrics возвращает весь список метрик.
func (s *Storage) ListMetrics() []model.Metric {
	s.RLock()
	defer s.RUnlock()

	list := make([]model.Metric, 0, len(s.metrics))
	for _, val := range s.metrics {
		list = append(list, *val)
	}

	return list
}

// UpdateMetrics не реализован для in-memory хранилища.
func (*Storage) UpdateMetrics(_ context.Context, _ []model.Metric) error {
	return nil
}

// Ping не реализован для in-memory хранилища.
func (*Storage) Ping(_ context.Context) error {
	return nil
}

// Shutdown коректно завершает работу с хранилищем.
func (s *Storage) Shutdown() {
	if s.fileStor != nil {
		s.saveToFile()
		s.fileStor.Close()
	}
}

// saveDump выполняет периодическую запись метрик в файл.
func (s *Storage) saveDump(ctx context.Context) {
	storeIntervalTimer := time.NewTimer(time.Second * time.Duration(s.cfg.StoreInterval))
	defer storeIntervalTimer.Stop()

	if s.cfg.StoreInterval == 0 {
		storeIntervalTimer.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-storeIntervalTimer.C:
			s.saveToFile()
		}
	}
}

// saveToFile записывает метрики в файл.
func (s *Storage) saveToFile() {
	s.RLock()
	defer s.RUnlock()

	if err := s.fileStor.SaveToFile(s.metrics); err != nil {
		slog.Error("ошибка сохранения метрик в файл", slog.Any("error", err))
	}
	s.fileStor.Close()
}
