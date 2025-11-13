package postgresql

import (
	"database/sql"
	"log/slog"

	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Storage struct {
	cfg *sconfig.Config
	db  *sql.DB
}

func New(cfg *sconfig.Config) (*Storage, error) {
	if cfg.DatabaseDSN == "" {
		return &Storage{cfg: cfg}, nil
	}

	db, err := sql.Open("pgx", cfg.DatabaseDSN)
	if err != nil {
		return nil, err
	}

	loadMigrations(db)

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

func loadMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		slog.Error(err.Error())
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://./migrations", "postgres", driver)
	if err != nil {
		slog.Error(err.Error())
		return
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error(err.Error())
	}
}
