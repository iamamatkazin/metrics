// Package repository предоставляет хранилища для сохранения метрик.
// Определяет интерфейс Storager и реализует несколько типов хранилищ:
//   - In-memory хранилище (memstorage): хранилище в памяти с опциональным файловым бэкапом
//   - PostgreSQL хранилище (postgresql): постоянное хранилище в базе данных
//   - Файловое хранилище (filestorage): JSON файловое хранилище
//
// Пакет использует паттерн декоратора для объединения хранилищ.
package repository

import (
	"context"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/repository/memstorage"
	"github.com/iamamatkazin/metrics.git/internal/repository/postgresql"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Storager определяет интерфейс для работы с метриками.
type Storager interface {
	GetMetric(ctx context.Context, id string) (*model.Metric, error)
	UpdateMetric(ctx context.Context, metric *model.Metric) error
	UpdateMetrics(ctx context.Context, metric []model.Metric) error
	ListMetrics() []model.Metric
	Ping(ctx context.Context) error
	Shutdown()
}

// Storage объединяет несколько хранилищ метрик.
type Storage struct {
	cfg      *server.Config
	memStor  *memstorage.Storage
	dbStor   *postgresql.Storage
	storages []Storager
}

// New создает новое хранилище с заданной конфигурацией.
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

// Shutdown завершает работу всех хранилищ.
func (s *Storage) Shutdown() {
	for _, storage := range s.storages {
		storage.Shutdown()
	}
}

// GetMetric возвращает метрику по её идентификатору.
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

// UpdateMetric создает или обновляет метрику.
func (s *Storage) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	for _, storage := range s.storages {
		if err := storage.UpdateMetric(ctx, metric); err != nil {
			return err
		}
	}

	return nil
}

// UpdateMetrics создает или обновляет массив метрик.
func (s *Storage) UpdateMetrics(ctx context.Context, metrics []model.Metric) error {
	for _, storage := range s.storages {
		if err := storage.UpdateMetrics(ctx, metrics); err != nil {
			return err
		}
	}

	return nil
}

// ListMetrics возвращает список всех метрик.
func (s *Storage) ListMetrics() []model.Metric {
	for _, storage := range s.storages {
		if list := storage.ListMetrics(); list != nil {
			return list
		}
	}

	return nil
}

// Ping проверяет доступность хранилищ.
func (s *Storage) Ping(ctx context.Context) error {
	for _, storage := range s.storages {
		if err := storage.Ping(ctx); err != nil {
			return err
		}
	}

	return nil
}
