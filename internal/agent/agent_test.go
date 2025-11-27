package agent

import (
	"context"
	"testing"

	"github.com/iamamatkazin/metrics.git/pkg/config/agent"
)

// mockClient реализует интерфейс Post для тестирования
type mockClient struct {
	lastURL         string
	lastContentType string
	err             error
}

func (m *mockClient) Post(ctx context.Context, url, contentType string, data any) error {
	m.lastURL = url
	m.lastContentType = contentType
	return m.err
}

func TestNew(t *testing.T) {
	type args struct {
		cfg *agent.Config
	}
	tests := []struct {
		name string
		args args
		want *Agent
	}{
		{
			name: "Test 1",
			args: args{cfg: &agent.Config{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			New(tt.args.cfg)
		})
	}
}

// func TestAgent_sendMetric(t *testing.T) {
// 	type args struct {
// 		urlBase string
// 		key     string
// 		name    string
// 		value   float64
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		mockErr error
// 	}{
// 		{
// 			name:    "Test_err_false",
// 			args:    args{},
// 			mockErr: nil,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock := &mockClient{
// 				err: tt.mockErr,
// 			}
// 			a := &Agent{
// 				client: mock,
// 			}

// 			a.sendMetric(tt.args.urlBase, tt.args.key, tt.args.name, tt.args.value)
// 		})
// 	}
// }

// func TestAgent_sendMetrics(t *testing.T) {
// 	type args struct {
// 		ctx context.Context
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		mockErr error
// 		wantErr bool
// 	}{
// 		{
// 			name:    "Test_err_false",
// 			args:    args{},
// 			mockErr: nil,
// 			wantErr: false,
// 		},
// 		{
// 			name:    "Test_err_true",
// 			args:    args{},
// 			mockErr: fmt.Errorf("error"),
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			mock := &mockClient{
// 				err: tt.mockErr,
// 			}

// 			metrics := make(map[string]map[string]float64)
// 			metrics[model.Gauge] = make(map[string]float64)
// 			metrics[model.Counter] = make(map[string]float64)
// 			metrics[model.Gauge]["test"] = 1.

// 			a := &Agent{
// 				client:  mock,
// 				metrics: metrics,
// 				cfg:     &agent.Config{},
// 			}
// 			a.sendMetrics()
// 		})
// 	}
// }

func TestAgent_reportMetrics(t *testing.T) {
	tests := []struct {
		name    string
		mockErr error
	}{
		{
			name:    "Test_err_false",
			mockErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// mock := &mockClient{
			// 	err: tt.mockErr,
			// }

			// 			metrics := make(map[string]map[string]float64)
			// 			metrics[model.Gauge] = make(map[string]float64)
			// 			metrics[model.Counter] = make(map[string]float64)
			// 			metrics[model.Gauge]["test"] = 1.

			// 			a := &Agent{
			// 				client:  mock,
			// 				metrics: metrics,
			// 				cfg:     &agent.Config{},
			// 			}

			// a.reportMetrics()
		})
	}
}
