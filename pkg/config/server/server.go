package server

import (
	"flag"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Address string `env:"ADDRESS"`
}

func New() (*Config, error) {
	address := flag.String("a", "localhost:8080", "a address")
	flag.Parse()

	cfg := &Config{Address: *address}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
