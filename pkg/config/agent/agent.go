package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	Timeout        time.Duration
	PollInterval   int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	report := flag.Int("r", 10, "частота отправки метрик на сервер")
	pool := flag.Int("p", 2, "частота опроса метрик")
	key := flag.String("k", "", "ключ подписи данных")
	flag.Parse()

	cfg := &Config{
		Address:        *address,
		Key:            *key,
		Timeout:        time.Second * 10,
		PollInterval:   *pool,
		ReportInterval: *report,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
