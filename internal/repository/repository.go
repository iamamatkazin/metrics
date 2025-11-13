package repository

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/repository/filestorage"
	"github.com/iamamatkazin/metrics.git/internal/repository/memstorage"
	"github.com/iamamatkazin/metrics.git/internal/repository/postgresql"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Storager interface {
	GetMetric(id string) *model.Metric
	UpdateMetric(ctx context.Context, metric model.Metric) error
	ListMetrics() []model.Metric
	PingDB(ctx context.Context) error
	Shutdown(ctx context.Context)
}

type Storage struct {
	cfg      *server.Config
	memStor  *memstorage.Storage
	fileStor *filestorage.Storage
	dbStor   *postgresql.Storage
}

func New(ctx context.Context, cfg *server.Config) (*Storage, error) {
	fileStor, err := filestorage.New(cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("cfg.DatabaseDSN", cfg.DatabaseDSN)

	dbStor, err := postgresql.New(cfg)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]*model.Metric)
	if cfg.Restore {
		if err = fileStor.LoadDump(metrics); err != nil {
			return nil, err
		}
	}

	mem := memstorage.New(metrics)
	s := &Storage{
		cfg:      cfg,
		memStor:  mem,
		fileStor: fileStor,
		dbStor:   dbStor,
	}

	go func() {
		s.saveDump(ctx)
	}()

	return s, nil
}

func (s *Storage) Shutdown(ctx context.Context) {
	metrics := s.memStor.GetMetrics()

	if s.fileStor != nil {
		s.saveToFile(metrics)
		s.fileStor.Close()
	}

	if s.dbStor != nil {
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()

		s.dbStor.UpdateMetrics(ctx, metrics)
		s.dbStor.Close()
		slog.Info("UpdateMetrics")
	}
}

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
			metrics := s.memStor.GetMetrics()
			s.saveToFile(metrics)
			s.dbStor.UpdateMetrics(ctx, metrics)
		}
	}
}

func (s *Storage) saveToFile(metrics map[string]*model.Metric) {
	if err := s.fileStor.SaveToFile(metrics); err != nil {
		slog.Error("ошибка сохранения метрик в файл", slog.Any("error", err))
	}
}

func (s *Storage) GetMetric(id string) *model.Metric {
	return s.memStor.GetMetric(id)
}

func (s *Storage) UpdateMetric(ctx context.Context, metric model.Metric) error {
	val := s.memStor.UpdateMetric(metric)

	if s.cfg.StoreInterval == 0 {
		metrics := s.memStor.GetMetrics()
		s.saveToFile(metrics)
	}

	if err := s.dbStor.UpdateMetric(ctx, val); err != nil {
		return err
	}

	return nil
}

func (s *Storage) ListMetrics() []model.Metric {
	return s.memStor.ListMetrics()
}

func (s *Storage) PingDB(ctx context.Context) error {
	return s.dbStor.Ping(ctx)
}
