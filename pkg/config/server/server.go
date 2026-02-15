// Package server предоставляет структуры конфигурации для сервера метрик.
// Включает тип Config, который содержит все настройки сервера, включая
// сетевые адреса, пути к хранилищу, подключение к базе данных и настройки аудита.
package server

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config - структура конфигурации сервера.
type Config struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	FileAudit       string `env:"AUDIT_FILE"`
	URLAudit        string `env:"AUDIT_URL"`
	Timeout         time.Duration
	StoreInterval   int  `env:"STORE_INTERVAL"`
	Restore         bool `env:"RESTORE"`
}

// New создает новую конфигурацию сервера с настройками по умолчанию.
// Параметры могут быть переопределены через флаги командной строки и переменные окружения.
func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	interval := flag.Int("i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	path := flag.String("f", "./storage.json", "путь до файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", true, "булево значение (true/false), определяющее, следует ли загружать ранее сохранённые значения из указанного файла при старте сервера")
	database := flag.String("d", "", "строка с адресом подключения к БД") // host=localhost user=postgres password=postgres dbname=metrics sslmode=disable
	key := flag.String("k", "", "ключ подписи данных")
	fileAudit := flag.String("audit-file", "", "путь к файлу, в который сохраняются логи аудита")
	urlAudit := flag.String("audit-url", "", "полный URL, по которому отправляются логи аудита")

	flag.Parse()

	cfg := &Config{
		Address:         *address,
		StoreInterval:   *interval,
		FileStoragePath: *path,
		Key:             *key,
		Restore:         *restore,
		DatabaseDSN:     *database,
		FileAudit:       *fileAudit,
		URLAudit:        *urlAudit,
		Timeout:         time.Second * 10,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
