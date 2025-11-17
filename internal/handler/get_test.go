package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/config/server"
	"github.com/stretchr/testify/require"
)

func prepareMetric(h *Handler, mType, name string) {
	r := httptest.NewRequest(http.MethodPost, "/update/"+mType+"/"+name+"/100", nil)
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("type", mType)
	ctx.URLParams.Add("id", name)
	ctx.URLParams.Add("val", "100")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

	h.updateMetric(w, r)
}
func TestHandler_getMetric(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
		request    string
		mType      string
		val        string
		want       want
	}{
		{
			name:       "simple test #1",
			request:    "/value/counter/testCounter",
			metricName: "testCounter",
			mType:      model.Counter,
			val:        "200",
			want:       want{statusCode: 200, contentType: "text/plain; charset=utf-8"},
		},
		{
			name:       "simple test #2",
			request:    "/value/gauge/test",
			metricName: "test",
			mType:      model.Gauge,
			val:        "100",
			want:       want{statusCode: 200, contentType: "text/plain; charset=utf-8"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, _ := New(context.Background(), &server.Config{FileStoragePath: "./storage.json"})
			prepareMetric(h, model.Counter, "testCounter")
			prepareMetric(h, model.Gauge, "test")
			prepareMetric(h, model.Counter, "testCounter")

			r := httptest.NewRequest(http.MethodGet, tt.request, nil)
			w := httptest.NewRecorder()

			ctx := chi.NewRouteContext()
			ctx.URLParams.Add("type", tt.mType)
			ctx.URLParams.Add("id", tt.metricName)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, ctx))

			h.getMetric(w, r)
			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			got, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.val, string(got))
		})
	}
}

func TestHandler_getMetricJSON(t *testing.T) {
	tests := []struct {
		name       string
		metricName string
		request    string
		data       string
		value      float64
		want       want
	}{
		// {
		// 	name:       "simple test #1",
		// 	request:    "/value/",
		// 	data:       "{\"id\":\"testCounter\",\"type\":\"gauge\"}",
		// 	metricName: "testCounter",
		// 	value:      100,
		// 	want:       want{statusCode: 200, contentType: "application/json", body: "{\"id\":\"testCounter\",\"type\":\"gauge\",\"value\":100}"},
		// },
		// {
		// 	name:       "simple test #2",
		// 	request:    "/value/",
		// 	data:       "{\"id\":\"testCounter\",\"type\":\"gauge\"}",
		// 	metricName: "testCounter",
		// 	value:      100,
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, err := New(context.Background(), &server.Config{FileStoragePath: "./storage.json"})
			if err != nil {
				t.Fatalf("could not construct receiver type: %v", err)
			}

			h.storage.UpdateMetric(context.Background(), model.Metric{
				ID:    tt.metricName,
				MType: model.Gauge,
				Value: &tt.value,
			})

			r := httptest.NewRequest(http.MethodGet, tt.request, strings.NewReader(tt.data))
			w := httptest.NewRecorder()

			h.getMetricJSON(w, r)

			result := w.Result()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))

			got, err := io.ReadAll(result.Body)
			require.NoError(t, err)
			err = result.Body.Close()
			require.NoError(t, err)

			assert.Equal(t, tt.want.body, string(got))
		})
	}
}
