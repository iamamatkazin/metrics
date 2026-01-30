package repository

import (
	"context"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/repository/memstorage"
	"github.com/iamamatkazin/metrics.git/internal/repository/postgresql"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Storager - интерфейс работы с метриками для разных хранилищ.
type Storager interface {
	GetMetric(ctx context.Context, id string) (*model.Metric, error)
	UpdateMetric(ctx context.Context, metric *model.Metric) error
	UpdateMetrics(ctx context.Context, metric []model.Metric) error
	ListMetrics() []model.Metric
	Ping(ctx context.Context) error
	Shutdown()
}

// Storage - структура, которая хранит в себе перечень всех возможных хранилищ для метрик.
type Storage struct {
	cfg      *server.Config
	memStor  *memstorage.Storage
	dbStor   *postgresql.Storage
	storages []Storager
}

// New - конструктор для Storage.
func New(ctx context.Context, cfg *server.Config) (*Storage, error) {
	dbStor, err := postgresql.New(cfg)
	if err != nil {
		return nil, err
	}

	mem, err := memstorage.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	storages := make([]Storager, 0, 2)
	storages = append(storages, mem)
	storages = append(storages, dbStor)

	s := &Storage{
		cfg:      cfg,
		memStor:  mem,
		dbStor:   dbStor,
		storages: storages,
	}

	return s, nil
}

// Shutdown - завершение работы всех доступных хранилищ.
func (s *Storage) Shutdown() {
	for _, storage := range s.storages {
		storage.Shutdown()
	}
}

// GetMetric - получить метрику по ее идентификатору.
func (s *Storage) GetMetric(ctx context.Context, id string) (*model.Metric, error) {
	for _, storage := range s.storages {
		metric, err := storage.GetMetric(ctx, id)
		if err != nil {
			return nil, err
		}

		if metric != nil {
			return metric, nil
		}
	}

	return nil, nil
}

// UpdateMetric - создание метрики, если она есть, то изменение.
func (s *Storage) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	for _, storage := range s.storages {
		if err := storage.UpdateMetric(ctx, metric); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMetrics - создание массива метрик, если они есть, то изменение.
func (s *Storage) UpdateMetrics(ctx context.Context, metrics []model.Metric) error {
	for _, storage := range s.storages {
		if err := storage.UpdateMetrics(ctx, metrics); err != nil {
			return err
		}
	}

	return nil
}

// ListMetrics - получить список метрик.
func (s *Storage) ListMetrics() []model.Metric {
	for _, storage := range s.storages {
		if list := storage.ListMetrics(); list != nil {
			return list
		}
	}

	return nil
}

// Ping - проверить доступность базы данных.
func (s *Storage) Ping(ctx context.Context) error {
	for _, storage := range s.storages {
		if err := storage.Ping(ctx); err != nil {
			return err
		}
	}

	return nil
}
