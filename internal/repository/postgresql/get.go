package postgresql

import (
	"context"
	"database/sql"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (s *Storage) GetMetric(ctx context.Context, id string) (*model.Metric, error) {
	if s.db == nil {
		return nil, nil
	}

	var (
		metric model.Metric
		val    sql.NullFloat64
		delta  sql.NullInt64
	)

	query := `SELECT id, mtype, val, delta FROM metrics WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	err := row.Scan(&metric.ID, &metric.MType, &val, &delta)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if val.Valid {
		metric.Value = &val.Float64
	}

	if delta.Valid {
		d := int(delta.Int64)
		metric.Delta = &d
	}

	return &metric, nil
}
