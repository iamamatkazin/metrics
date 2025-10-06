package config

import "time"

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
	Port int
}

func New() *Config {
	return &Config{
		Client: Client{
			Timeout:        time.Second * 5,
			PollInterval:   time.Second * 2,
			ReportInterval: time.Second * 10,
		},
		Server: Server{
			Host: "localhost",
			Port: 8080,
		},
	}
}
