package config

import (
	"flag"
)

type Server struct {
	Host string
}

func NewServer() *Server {
	host := flag.String("a", "localhost:8080", "a host")
	flag.Parse()

	return &Server{
		Host: *host,
	}
}
