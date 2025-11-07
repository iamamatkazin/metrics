package repository

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Storager interface {
	GetMetric(id string) *model.Metric
	UpdateMetric(metric model.Metric)
	ListMetrics() []model.Metric
}

type MemStorage struct {
	metrics map[string]*model.Metric
	sync.RWMutex
	cfg  *server.Config
	file *os.File
}

func New(ctx context.Context, cfg *server.Config) (*MemStorage, error) {
	file, err := os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	s := &MemStorage{
		metrics: make(map[string]*model.Metric),
		cfg:     cfg,
		file:    file,
	}

	if cfg.Restore {
		s.loadDump()
	}

	go func() {
		s.saveDump(ctx)
	}()

	return s, nil
}

func (s *MemStorage) Close() {
	if s.file != nil {
		s.file.Close()
	}
}

func (s *MemStorage) saveDump(ctx context.Context) {
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
			s.Lock()
			if err := s.saveToFile(); err != nil {
				slog.Error("ошибка сохранения метрик в файл", slog.Any("error", err))
			}
			s.Unlock()
		}
	}
}

func (s *MemStorage) loadDump() {
	if err := json.NewDecoder(s.file).Decode(&s.metrics); err != nil {
		if err != io.EOF {
			slog.Error("ошибка чтения метрик из файла", slog.Any("error", err))
		}
	}
}

func (s *MemStorage) saveToFile() error {
	if err := s.file.Truncate(0); err != nil {
		return err
	}
	if _, err := s.file.Seek(0, 0); err != nil {
		return err
	}

	if err := json.NewEncoder(s.file).Encode(s.metrics); err != nil {
		return err
	}

	return nil
}

func (s *MemStorage) GetMetric(id string) *model.Metric {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.metrics[id]
	if !ok {
		return nil
	}

	return val
}

func (s *MemStorage) UpdateMetric(metric model.Metric) {
	s.Lock()
	defer s.Unlock()

	val, ok := s.metrics[metric.ID]
	if !ok {
		s.metrics[metric.ID] = &metric
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

	if s.cfg.StoreInterval == 0 {
		if err := s.saveToFile(); err != nil {
			slog.Error("ошибка сохранения метрик в файл", slog.Any("error", err))
		}
	}
}

func (s *MemStorage) ListMetrics() []model.Metric {
	s.RLock()
	defer s.RUnlock()

	list := make([]model.Metric, 0, len(s.metrics))
	for _, val := range s.metrics {
		list = append(list, *val)
	}

	return list
}
