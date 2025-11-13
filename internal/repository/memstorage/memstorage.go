package memstorage

import (
	"maps"
	"sync"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Storage struct {
	metrics map[string]*model.Metric
	sync.RWMutex
	cfg *server.Config
}

func New(metrics map[string]*model.Metric) *Storage {
	return &Storage{
		metrics: metrics,
	}
}

func (s *Storage) GetMetrics() map[string]*model.Metric {
	s.RLock()
	defer s.RUnlock()

	copy := make(map[string]*model.Metric, len(s.metrics))
	maps.Copy(copy, s.metrics)

	return copy
}

func (s *Storage) GetMetric(id string) *model.Metric {
	s.RLock()
	defer s.RUnlock()

	val, ok := s.metrics[id]
	if !ok {
		return nil
	}

	return val
}

func (s *Storage) UpdateMetric(metric model.Metric) *model.Metric {
	s.Lock()
	defer s.Unlock()

	val, ok := s.metrics[metric.ID]
	if !ok {
		s.metrics[metric.ID] = &metric
		val = &metric
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

	return val
}

func (s *Storage) ListMetrics() []model.Metric {
	s.RLock()
	defer s.RUnlock()

	list := make([]model.Metric, 0, len(s.metrics))
	for _, val := range s.metrics {
		list = append(list, *val)
	}

	return list
}
