package client

import (
	"context"
	"errors"
	"testing"

	gproto "github.com/iamamatkazin/metrics.git/internal/grpc/proto"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// mockMetricsClient реализует интерфейс proto.MetricsClient для мокирования.
type mockMetricsClient struct {
	gproto.UnimplementedMetricsServer
	updateMetricsFunc func(ctx context.Context, in *gproto.UpdateMetricsRequest, opts ...grpc.CallOption) (*gproto.UpdateMetricsResponse, error)
}

func (m *mockMetricsClient) UpdateMetrics(ctx context.Context, in *gproto.UpdateMetricsRequest, opts ...grpc.CallOption) (*gproto.UpdateMetricsResponse, error) {
	if m.updateMetricsFunc != nil {
		return m.updateMetricsFunc(ctx, in, opts...)
	}
	return &gproto.UpdateMetricsResponse{}, nil
}

// TestUpdateMetrics_Success тестирует успешный сценарий обновления метрик.
func TestUpdateMetrics_Success(t *testing.T) {
	// Создаем мок-клиент с успешным ответом
	mockClient := &mockMetricsClient{
		updateMetricsFunc: func(ctx context.Context, in *gproto.UpdateMetricsRequest, opts ...grpc.CallOption) (*gproto.UpdateMetricsResponse, error) {
			// Проверяем, что контекст содержит метаданные
			md, ok := metadata.FromOutgoingContext(ctx)
			assert.True(t, ok, "ожидались метаданные в контексте")
			assert.NotNil(t, md, "метаданные не должны быть nil")

			// Проверяем, что запрос содержит метрики
			assert.NotNil(t, in, "запрос не должен быть nil")

			return &gproto.UpdateMetricsResponse{}, nil
		},
	}

	// Создаем тестируемый объект Metrics с моком
	metrics := &Metrics{
		client: mockClient,
	}

	// Создаем тестовые данные
	testMetrics := []model.Metric{
		{
			ID:    "testGauge",
			MType: model.Gauge,
			Value: func() *float64 { v := 1.5; return &v }(),
		},
		{
			ID:    "testCounter",
			MType: model.Counter,
			Delta: func() *int { v := 100; return &v }(),
		},
	}

	// Вызываем тестируемую функцию
	err := metrics.UpdateMetrics(context.Background(), testMetrics)

	// Проверяем результат
	assert.NoError(t, err, "ожидалось отсутствие ошибки")
}

// TestUpdateMetrics_Error тестирует сценарий с ошибкой при обновлении метрик.
func TestUpdateMetrics_Error(t *testing.T) {
	// Создаем мок-клиент, который возвращает ошибку
	expectedErr := errors.New("gRPC error: connection refused")
	mockClient := &mockMetricsClient{
		updateMetricsFunc: func(ctx context.Context, in *gproto.UpdateMetricsRequest, opts ...grpc.CallOption) (*gproto.UpdateMetricsResponse, error) {
			return nil, expectedErr
		},
	}

	// Создаем тестируемый объект Metrics с моком
	metrics := &Metrics{
		client: mockClient,
	}

	// Создаем тестовые данные
	testMetrics := []model.Metric{
		{
			ID:    "testGauge",
			MType: model.Gauge,
			Value: func() *float64 { v := 1.5; return &v }(),
		},
	}

	// Вызываем тестируемую функцию
	err := metrics.UpdateMetrics(context.Background(), testMetrics)

	// Проверяем результат
	assert.Equal(t, expectedErr, err, "gRPC error: connection refused")
}

// TestUpdateMetrics_EmptyList тестирует сценарий с пустым списком метрик.
func TestUpdateMetrics_EmptyList(t *testing.T) {
	// Создаем мок-клиент
	mockClient := &mockMetricsClient{
		updateMetricsFunc: func(ctx context.Context, in *gproto.UpdateMetricsRequest, opts ...grpc.CallOption) (*gproto.UpdateMetricsResponse, error) {
			// Проверяем, что список метрик пустой
			assert.NotNil(t, in, "запрос не должен быть nil")
			return &gproto.UpdateMetricsResponse{}, nil
		},
	}

	// Создаем тестируемый объект Metrics с моком
	metrics := &Metrics{
		client: mockClient,
	}

	// Вызываем тестируемую функцию с пустым списком
	err := metrics.UpdateMetrics(context.Background(), []model.Metric{})

	// Проверяем результат
	assert.NoError(t, err, "ожидалось отсутствие ошибки с пустым списком")
}
