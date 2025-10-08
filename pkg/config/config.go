package config

import (
	"flag"
	"time"
)

type Config struct {
	Client Client
	Server Server
}

type Client struct {
	Timeout        time.Duration
	PollInterval   time.Duration
	ReportInterval time.Duration
}

type Server struct {
	Host string
}

func New() *Config {
	host := flag.String("a", "localhost:8080", "a host")
	report := flag.Int("r", 10, "a report")
	pool := flag.Int("p", 2, "a pool")

	flag.Parse()

	c := &Config{
		Client: Client{
			Timeout:        time.Second * 5,
			PollInterval:   time.Second * time.Duration(*pool),
			ReportInterval: time.Second * time.Duration(*report),
		},
		Server: Server{
			Host: *host,
		},
	}

	return c
}
