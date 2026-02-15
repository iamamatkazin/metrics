// Package postgresql предоставляет PostgreSQL хранилище для сохранения метрик.
// Реализует интерфейс Storager с использованием базы данных PostgreSQL,
// поддерживая CRUD операции для метрик с логикой повторных попыток.
package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/iamamatkazin/metrics.git/internal/model"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Storage - структура, реализующая интерфейс работы с метриками для хранилища в PostgreSQL.
type Storage struct {
	cfg *sconfig.Config
	db  *sql.DB
}

// New создает новое подключение к базе данных PostgreSQL.
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

// loadMigrations загружает миграции в базу данных.
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

// isRetryablePgError проверяет, можно ли повторить запрос после ошибки.
func isRetryablePgError(err error) bool {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		// Неизвестная ошибка — не повторяем
		return false
	}

	switch pgErr.Code {
	case pgerrcode.ConnectionException,
		pgerrcode.ConnectionDoesNotExist,
		pgerrcode.ConnectionFailure,
		pgerrcode.SQLClientUnableToEstablishSQLConnection,
		pgerrcode.SQLServerRejectedEstablishmentOfSQLConnection,
		pgerrcode.TransactionResolutionUnknown,
		pgerrcode.ProtocolViolation:
		return true
	default:
		return false
	}
}

// ListMetrics - метод не реализован для PostgreSQL.
func (*Storage) ListMetrics() []model.Metric { return nil }

// PingDB - метод не реализован для PostgreSQL.
func (*Storage) PingDB(context.Context) error { return nil }

// Shutdown коректно завершает работу с хранилищем.
func (s *Storage) Shutdown() {
	if s.db != nil {
		s.db.Close()
	}
}
