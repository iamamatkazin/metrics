package postgresql

import (
	"context"
	"database/sql"
	"time"

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
		err    error
	)

	metric.ID = id
	query := `SELECT mtype, val, delta FROM metrics WHERE id = $1`

	metric.MType, val, delta, err = s.retryableGet(ctx, query, id)
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

func (s *Storage) retryableGet(ctx context.Context, query, id string) (mType string, val sql.NullFloat64, delta sql.NullInt64, err error) {
	timerRetryable := time.NewTimer(0)
	count := 0

	for {
		select {
		case <-ctx.Done():
		case <-timerRetryable.C:
			row := s.db.QueryRowContext(ctx, query, id)

			err = row.Scan(&mType, &val, &delta)
			if err != nil {
				if !isRetryablePgError(err) || count > 3 {
					return "", sql.NullFloat64{}, sql.NullInt64{}, err
				}

				timerRetryable.Reset(time.Duration(2*count+1) * time.Second)
				count++
			}

			return mType, val, delta, nil
		}
	}
}
