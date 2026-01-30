package postgresql

import (
	"context"
	"database/sql"
	"time"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

// UpdateMetric - создание метрики, если она есть, то изменение.
func (s *Storage) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	if s.db == nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := s.updateMetric(ctx, tx, *metric); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// UpdateMetrics - создание массива метрик, если они есть, то изменение.
func (s *Storage) UpdateMetrics(ctx context.Context, metrics []model.Metric) error {
	if s.db == nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for i := range metrics {
		if err := s.updateMetric(ctx, tx, metrics[i]); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

// updateMetric - внутренний метод создания метрики, если она есть, то изменение.
func (s *Storage) updateMetric(ctx context.Context, tx *sql.Tx, metric model.Metric) error {
	query := `
		INSERT INTO metrics (id, mtype, val, delta)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id) DO UPDATE SET
			mtype = EXCLUDED.mtype,
			val = EXCLUDED.val,
			delta = EXCLUDED.delta
	`

	if metric.MType == model.Counter {
		if metric.Value != nil {
			delta := int(*metric.Value)
			metric.Delta = &delta
			metric.Value = nil
		}

		query = `
			INSERT INTO metrics (id, mtype, val, delta)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				mtype = EXCLUDED.mtype,
				val = EXCLUDED.val,
				delta = metrics.delta + EXCLUDED.delta
		`
	}

	if err := retryableUpdateTx(ctx, tx, query, metric); err != nil {
		return err
	}

	return nil
}

// retryableUpdateTx - проверка возможности повторного запроса в базу данных.
func retryableUpdateTx(ctx context.Context, tx *sql.Tx, query string, metric model.Metric) error {
	timerRetryable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetryable.C:
			_, err := tx.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value, metric.Delta)
			if err != nil {
				if !isRetryablePgError(err) || count > 3 {
					return err
				}

				timerRetryable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return nil
		}
	}
}
