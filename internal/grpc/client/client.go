package client

import (
	"fmt"

	"github.com/iamamatkazin/metrics.git/internal/grpc/proto"
	"github.com/iamamatkazin/metrics.git/pkg/config/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Metrics реализует клиента grpc.
type Metrics struct {
	conn   *grpc.ClientConn
	client proto.MetricsClient
}

func New(cfg *agent.Config) (*Metrics, error) {
	// Устанавливаем соединение с сервером
	conn, err := grpc.NewClient(fmt.Sprintf(":%d", cfg.GrpcPort), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := proto.NewMetricsClient(conn)

	return &Metrics{conn: conn, client: client}, nil
}

func (c *Metrics) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
