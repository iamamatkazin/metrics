// Package server предоставляет структуры конфигурации для сервера метрик.
// Включает тип Config, который содержит все настройки сервера, включая
// сетевые адреса, пути к хранилищу, подключение к базе данных и настройки аудита.
package server

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config - структура конфигурации сервера.
type Config struct {
	Address         string `json:"address" env:"ADDRESS"`
	FileStoragePath string `json:"file_storage_path" env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `json:"database_dsn" env:"DATABASE_DSN"`
	Key             string `json:"key" env:"KEY"`
	FileAudit       string `json:"audit_file" env:"AUDIT_FILE"`
	URLAudit        string `json:"audit_url" env:"AUDIT_URL"`
	CryptoKey       string `json:"crypto_key" env:"CRYPTO_KEY"`
	FileConfig      string `env:"CONFIG"`
	TrustedSubnet   string `json:"trusted_subnet" env:"TRUSTED_SUBNET"`
	Timeout         time.Duration
	StoreInterval   int  `json:"store_interval" env:"STORE_INTERVAL"`
	GrpcPort        int  `json:"grpc_port" env:"GRPC_PORT"`
	Restore         bool `json:"restore" env:"RESTORE"`
}

// New создает новую конфигурацию сервера с настройками по умолчанию.
// Параметры могут быть переопределены через флаги командной строки и переменные окружения.
func New() (*Config, error) {
	var cfg Config

	config := flag.String("config", "", "конфигурация сервера с помощью файла в формате JSON")
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	interval := flag.Int("i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	path := flag.String("f", "./storage.json", "путь до файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", true, "булево значение (true/false), определяющее, следует ли загружать ранее сохранённые значения из указанного файла при старте сервера")
	database := flag.String("d", "", "строка с адресом подключения к БД")
	key := flag.String("k", "", "ключ подписи данных")
	fileAudit := flag.String("audit-file", "", "путь к файлу, в который сохраняются логи аудита")
	urlAudit := flag.String("audit-url", "", "полный URL, по которому отправляются логи аудита")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с приватным ключом")
	trustedSubnet := flag.String("trusted_subnet", "192.168.1.0/24", "содержит строковое представление бесклассовой адресации (CIDR).")
	grpcPort := flag.Int("g", 3200, "порт grpc сервера")

	flag.Parse()

	if *config != "" {
		data, err := os.ReadFile(*config)
		if err != nil {
			slog.Error("ошибка чтения файла конфигурации:", slog.Any("error", err))
		} else {
			err = json.Unmarshal(data, &cfg)
			if err != nil {
				slog.Error("ошибка парсинга файла конфигурации:", slog.Any("error", err))
			}
		}
	}

	cfg.Address = *address
	cfg.StoreInterval = *interval
	cfg.FileStoragePath = *path
	cfg.Key = *key
	cfg.Restore = *restore
	cfg.Timeout = time.Second * 10
	cfg.TrustedSubnet = *trustedSubnet
	cfg.GrpcPort = *grpcPort

	if *database != "" {
		cfg.DatabaseDSN = *database
	}

	if *fileAudit != "" {
		cfg.FileAudit = *fileAudit
	}

	if *urlAudit != "" {
		cfg.URLAudit = *urlAudit
	}

	if *cryptoKey != "" {
		cfg.CryptoKey = *cryptoKey
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
