package agent

import (
	"flag"
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address        string `env:"ADDRESS"`
	Timeout        time.Duration
	PollInterval   int `env:"POLL_INTERVAL"`
	ReportInterval int `env:"REPORT_INTERVAL"`
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "a address")
	report := flag.Int("r", 10, "a report")
	pool := flag.Int("p", 2, "a pool")
	flag.Parse()

	cfg := &Config{
		Address:        *address,
		Timeout:        time.Second * 5,
		PollInterval:   *pool,
		ReportInterval: *report,
	}

	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
