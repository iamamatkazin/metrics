package handler

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

type want struct {
	contentType string
	statusCode  int
}

func TestHandler_updateMetric(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
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

			h := New()
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
