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

// Storage - структура, реализующая интерфейс работы с метриками для хранилища в Posgresql.
type Storage struct {
	cfg *sconfig.Config
	db  *sql.DB
}

// New - конструктор для Storage.
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

// loadMigrations - загрузка миграций в базу.
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

// isRetryablePgError - проверка возможности повторного запроса в базу данных.
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

// UpdateMetrics - заглушка для метода.
func (s *Storage) ListMetrics() []model.Metric { return nil }

// Ping - заглушка для метода.
func (s *Storage) PingDB(ctx context.Context) error { return nil }

// Shutdown - завершить работу с хранилищем.
func (s *Storage) Shutdown() {
	if s.db != nil {
		s.db.Close()
	}
}
