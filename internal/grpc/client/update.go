package client

import (
	"context"

	gproto "github.com/iamamatkazin/metrics.git/internal/grpc/proto"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/http"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

func (c *Metrics) UpdateMetrics(ctx context.Context, list []model.Metric) error {
	md := metadata.New(map[string]string{"x-real-ip": http.GetIPAdress()})
	ctx = metadata.NewOutgoingContext(ctx, md)

	metrics := make([]*gproto.Metric, 0, len(list))

	for i := range list {
		metrics = append(metrics, fillProtoMetric(list[i]))
	}

	_, err := c.client.UpdateMetrics(ctx, gproto.UpdateMetricsRequest_builder{
		Metrics: metrics,
	}.Build())
	if err != nil {
		return err
	}

	return nil
}

func fillProtoMetric(metric model.Metric) *gproto.Metric {
	var mtype *gproto.Metric_MType

	if metric.MType == model.Gauge {
		mtype = gproto.Metric_GAUGE.Enum()
	} else {
		mtype = gproto.Metric_COUNTER.Enum()
	}

	var delta *int64
	if metric.Delta != nil {
		delta = proto.Int64(int64(*metric.Delta))
	}

	var value *float64
	if metric.Value != nil {
		value = proto.Float64(float64(*metric.Value))
	}

	return gproto.Metric_builder{
		Id:    proto.String(metric.ID),
		Type:  mtype,
		Delta: delta,
		Value: value,
	}.Build()
}
