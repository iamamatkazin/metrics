package postgresql

import (
	"context"
	"database/sql"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (s *Storage) UpdateMetric(ctx context.Context, metric *model.Metric) error {
	if s.db == nil {
		return nil
	}

	query := `
			INSERT INTO metrics (id, mtype, val, delta)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				mtype = EXCLUDED.mtype,
				val = EXCLUDED.val,
				delta = EXCLUDED.delta
		`

	if _, err := s.db.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value, metric.Delta); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateMetrics(ctx context.Context, metric map[string]*model.Metric) error {
	if s.db == nil {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	for key := range metric {
		if err := s.updateMetric(ctx, tx, metric[key]); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *Storage) updateMetric(ctx context.Context, tx *sql.Tx, metric *model.Metric) error {
	query := `
			INSERT INTO metrics (id, mtype, val, delta)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (id) DO UPDATE SET
				mtype = EXCLUDED.mtype,
				val = EXCLUDED.val,
				delta = EXCLUDED.delta
		`
	if _, err := tx.ExecContext(ctx, query, metric.ID, metric.MType, metric.Value, metric.Delta); err != nil {
		return err
	}

	return nil
}
