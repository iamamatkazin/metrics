package server

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address         string `env:"ADDRESS"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
	Key             string `env:"KEY"`
	StoreInterval   int    `env:"STORE_INTERVAL"`
	Restore         bool   `env:"RESTORE"`
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	interval := flag.Int("i", 300, "интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск")
	path := flag.String("f", "./storage.json", "путь до файла, куда сохраняются текущие значения")
	restore := flag.Bool("r", true, "булево значение (true/false), определяющее, следует ли загружать ранее сохранённые значения из указанного файла при старте сервера")
	database := flag.String("d", "", "строка с адресом подключения к БД") // host=localhost user=postgres password=postgres dbname=metrics sslmode=disable
	key := flag.String("k", "", "ключ подписи данных")
	flag.Parse()

	cfg := &Config{
		Address:         *address,
		StoreInterval:   *interval,
		FileStoragePath: *path,
		Key:             *key,
		Restore:         *restore,
		DatabaseDSN:     *database,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
