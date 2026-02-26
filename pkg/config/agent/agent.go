// Package agent предоставляет структуры конфигурации для агента метрик.
// Включает тип Config, который содержит все настройки агента, включая
// сетевые адреса, интервалы и таймауты.
package agent

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
)

// Config - структура конфигурации агента.
type Config struct {
	Address        string `json:"address" env:"ADDRESS"`
	Key            string `json:"key" env:"KEY"`
	CryptoKey      string `json:"crypto_key" env:"CRYPTO_KEY"`
	FileConfig     string `env:"CONFIG"`
	Timeout        time.Duration
	PollInterval   int `json:"poll_interval" env:"POLL_INTERVAL"`
	ReportInterval int `json:"report_interval" env:"REPORT_INTERVAL"`
	RateLimit      int `json:"rate_limit" env:"RATE_LIMIT"`
}

// New создает новую конфигурацию агента с настройками по умолчанию.
// Параметры могут быть переопределены через флаги командной строки и переменные окружения.
func New() (*Config, error) {
	var cfg Config

	config := flag.String("config", "", "конфигурация агента с помощью файла в формате JSON")
	address := flag.String("a", "localhost:8080", "адрес эндпоинта HTTP-сервера")
	report := flag.Int("r", 1, "частота отправки метрик на сервер")
	pool := flag.Int("p", 1, "частота опроса метрик")
	key := flag.String("k", "", "ключ подписи данных")
	rateLimit := flag.Int("l", 10, "количество одновременно исходящих запросов на сервер")
	cryptoKey := flag.String("crypto-key", "", "путь до файла с публичным ключом")

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
	cfg.Timeout = time.Second * 10
	cfg.PollInterval = *pool
	cfg.ReportInterval = *report
	cfg.RateLimit = *rateLimit

	if *key != "" {
		cfg.Key = *key
	}

	if *cryptoKey != "" {
		cfg.CryptoKey = *cryptoKey
	}

	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
