package repository

import (
	"sync"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

type Storager interface {
	GetMetric(id string) *model.Metric
	UpdateMetric(metric model.Metric)
	ListMetrics() []model.Metric
}

type MemStorage struct {
	metrics map[string]*model.Metric
	sync.RWMutex
}

func New() *MemStorage {
	return &MemStorage{
		metrics: make(map[string]*model.Metric),
	}
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
			delta := *val.Delta + *metric.Delta
			val.Delta = &delta
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
