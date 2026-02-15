package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
	"github.com/stretchr/testify/require"
)

type want struct {
	contentType string
	body        string
	statusCode  int
}

func TestHandler_updateMetric(t *testing.T) {
	tests := []struct {
		name    string
		request string
		mType   string
		val     string
		want    want
	}{
		{
			name:    "simple test #1",
			request: "/update/unknown/testCounter/100",
			mType:   "unknown",
			val:     "100",
			want:    want{statusCode: 400, contentType: "text/plain; charset=utf-8"},
		},
		{
			name:    "simple test #1",
			request: "/update/gauge/testCounter/1A0",
			mType:   "gauge",
			val:     "1A0",
			want:    want{statusCode: 400, contentType: "text/plain; charset=utf-8"},
		},
		{
			name:    "simple test #5",
			request: "/update/gauge/testCounter/100",
			mType:   "gauge",
			val:     "100",
			want:    want{statusCode: 200, contentType: "text/plain; charset=utf-8"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("type", tt.mType)
			ctx.URLParams.Add("id", "testCounter")
			ctx.URLParams.Add("val", tt.val)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			h, _ := New(context.Background(), &server.Config{FileStoragePath: "./storage.json"})
			h.updateMetric(w, r)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			_, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)
		})
	}
}

func Test_checkSign(t *testing.T) {
	type args struct {
		body any
		key  string
		hash string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test 1",
			args:    args{key: ""},
			wantErr: false,
		},
		{
			name:    "Test 2",
			args:    args{key: "key", body: make(chan int)},
			wantErr: true,
		},
		{
			name: "Test 3",
			args: args{
				key:  "key",
				body: "src",
				hash: "e2ba995098bc84e91b6cf607d3b6ef7a31f97ae38dfc7221f38e927bf2b49152",
			},
			wantErr: false,
		},
		{
			name: "Test 4",
			args: args{
				key:  "key",
				body: "src_bad",
				hash: "e2ba995098bc84e91b6cf607d3b6ef7a31f97ae38dfc7221f38e927bf2b49152",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, "/update/gauge/testCounter/1A0", nil)
			r.Header.Add("HashSHA256", tt.args.hash)

			if err := checkSign(r, tt.args.key, tt.args.body); (err != nil) != tt.wantErr {
				t.Errorf("checkSign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func ExampleHandler_updateMetricJSON() {
	var metric model.Metric

	// Передаем в теле запроса структуру вида:
	// {"id": "Alloc", "type": "gauge", "value": 12.4}
	err := json.NewDecoder(req.Body).Decode(&metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Валидируем переданный тип метрики, он должен принимать одно из двух значений:
	// Counter = "counter" или Gauge   = "gauge"
	if err = metric.ValidateJSON(); err != nil {
		// Обрабатываем ошибку
	}

	// Обновляем метрику в системе
	err = hand.storage.UpdateMetric(req.Context(), &metric)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Отправляем сообщение в систему аудита
	if message := hand.getMessage([]model.Metric{metric}, req.RemoteAddr); message != nil {
		hand.audit.Send(*message)
	}

	// Возвращаем ответ клиенту
}

func ExampleHandler_updatesMetricJSON() {
	var metrics []model.Metric

	// Передаем в теле запроса структуру вида:
	// [{"id": "Alloc", "type": "gauge", "value": 12.4}]
	err := json.NewDecoder(req.Body).Decode(&metrics)
	if err != nil {
		// Обрабатываем ошибку
	}

	// проверяем подпись
	if err = checkSign(req, hand.cfg.Key, metrics); err != nil {
		// Обрабатываем ошибку
	}

	// Обновляем метрики в системе
	err = hand.storage.UpdateMetrics(req.Context(), metrics)
	if err != nil {
		// Обрабатываем ошибку
	}

	// Отправляем сообщение в систему аудита
	if message := hand.getMessage(metrics, req.RemoteAddr); message != nil {
		hand.audit.Send(*message)
	}

	// Возвращаем ответ клиенту
}
