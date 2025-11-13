package postgresql

import (
	"context"
	"time"
)

func (s *Storage) Ping(ctx context.Context) error {
	if s.db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return err
	}

	return nil
}
