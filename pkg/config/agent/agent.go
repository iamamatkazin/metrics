// Package agent предоставляет структуры конфигурации для агента метрик.
// Включает тип Config, который содержит все настройки агента, включая
// сетевые адреса, интервалы и таймауты.
package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config - структура конфигурации агента.
type Config struct {
	Address        string `env:"ADDRESS"`
	Key            string `env:"KEY"`
	CryptoKey      string `env:"CRYPTO_KEY"`
	Timeout        time.Duration
	PollInterval   int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
	RateLimit      int `env:"RATE_LIMIT"`
}

// New создает новую конфигурацию агента с настройками по умолчанию.
// Параметры могут быть переопределены через флаги командной строки и переменные окружения.
func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	report := flag.Int("r", 1, "частота отправки метрик на сервер")
	pool := flag.Int("p", 1, "частота опроса метрик")
	key := flag.String("k", "", "ключ подписи данных")
	rateLimit := flag.Int("l", 10, "количество одновременно исходящих запросов на сервер")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с публичным ключом")
	flag.Parse()

	cfg := &Config{
		Address:        *address,
		Key:            *key,
		Timeout:        time.Second * 10,
		PollInterval:   *pool,
		ReportInterval: *report,
		RateLimit:      *rateLimit,
		CryptoKey:      *cryptoKey,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
