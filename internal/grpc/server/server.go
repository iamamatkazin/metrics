package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/iamamatkazin/metrics.git/internal/grpc/proto"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/repository"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Metrics поддерживает все необходимые методы сервера.
type Metrics struct {
	proto.UnimplementedMetricsServer
	server  *grpc.Server
	cfg     *sconfig.Config
	storage repository.Storager
}

func New(storage *repository.Storage, cfg *sconfig.Config) *Metrics {
	ms := Metrics{storage: storage, cfg: cfg}

	ms.server = grpc.NewServer(grpc.UnaryInterceptor(ms.unaryInterceptor))
	proto.RegisterMetricsServer(ms.server, &ms)

	return &ms
}

func (s *Metrics) Run(cfg *sconfig.Config) error {
	slog.Info("Запуск GRPC сервера")

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcPort))
	if err != nil {
		return err
	}

	if err := s.server.Serve(listen); err != nil {
		return err
	}

	return nil
}

func (s *Metrics) UpdateMetrics(ctx context.Context, in *proto.UpdateMetricsRequest) (*proto.UpdateMetricsResponse, error) {
	protoMetrics := in.GetMetrics()
	metrics := make([]model.Metric, 0, len(protoMetrics))

	for _, item := range protoMetrics {
		delta := int(item.GetDelta())
		value := item.GetValue()
		metrics = append(metrics, model.Metric{
			ID:    item.GetId(),
			MType: item.GetType().Enum().String(),
			Delta: &delta,
			Value: &value,
		})
	}

	if err := s.storage.UpdateMetrics(ctx, metrics); err != nil {
		return nil, status.Errorf(codes.Internal, "ошибка выполнения функции UpdateMetrics: %v", err)
	}

	return &proto.UpdateMetricsResponse{}, nil
}

func (s *Metrics) unaryInterceptor(ctx context.Context, req interface{}, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var ipAdress string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get("x-real-ip")
		if len(values) > 0 {
			ipAdress = values[0]
		}
	}

	if ipAdress != s.cfg.TrustedSubnet {
		return nil, status.Error(codes.PermissionDenied, "клиент не входит в доверенную сеть")
	}

	return handler(ctx, req)
}
