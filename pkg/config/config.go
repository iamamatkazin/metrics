package config

import (
	"flag"
	"time"
)

type Client struct {
	ServerAddress  string
	Timeout        time.Duration
	PollInterval   time.Duration
	ReportInterval time.Duration
}

func NewClient() *Client {
	host := flag.String("a", "localhost:8080", "a host")
	report := flag.Int("r", 10, "a report")
	pool := flag.Int("p", 2, "a pool")

	flag.Parse()

	return &Client{
		ServerAddress:  *host,
		Timeout:        time.Second * 5,
		PollInterval:   time.Second * time.Duration(*pool),
		ReportInterval: time.Second * time.Duration(*report),
	}
}
