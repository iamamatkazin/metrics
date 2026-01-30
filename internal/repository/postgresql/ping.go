package postgresql

import (
	"context"
	"time"
)

// Ping проверяет доступность подключения к базе данных PostgreSQL.
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
