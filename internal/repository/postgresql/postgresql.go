package postgresql

import (
	"database/sql"

	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	cfg *sconfig.Config
	db  *sql.DB
}

func New(cfg *sconfig.Config) (*Storage, error) {
	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	return &Storage{
		cfg: cfg,
		db:  db,
	}, nil
}

func (s *Storage) Close() {
	if s.db != nil {
		s.db.Close()
	}
}
