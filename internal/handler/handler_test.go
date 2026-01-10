package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/repository"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type errorWriter struct {
	header http.Header
}

func (e *errorWriter) Header() http.Header {
	if e.header == nil {
		e.header = make(http.Header)
	}
	return e.header
}

func (e *errorWriter) Write([]byte) (int, error) {
	return 0, fmt.Errorf("error")
}

func (e *errorWriter) WriteHeader(statusCode int) {
	// заглушка
}

func Test_writeHTML(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
		html   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				status: http.StatusBadRequest,
				w:      httptest.NewRecorder(),
			},
		},
		{
			name: "Test err",
			args: args{
				status: http.StatusBadRequest,
				w:      &errorWriter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeHTML(tt.args.w, tt.args.status, tt.args.html)
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		ctx context.Context
		cfg *sconfig.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Handler
		wantErr bool
	}{
		{
			name:    "Test 1",
			args:    args{cfg: &sconfig.Config{}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.ctx, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeJSON(t *testing.T) {
	type args struct {
		w      http.ResponseWriter
		status int
		body   []byte
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				status: http.StatusBadRequest,
				w:      httptest.NewRecorder(),
			},
		},
		{
			name: "Test err",
			args: args{
				status: http.StatusBadRequest,
				w:      &errorWriter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeJSON(tt.args.w, tt.args.status, tt.args.body)
		})
	}
}

func Test_writeText(t *testing.T) {
	type args struct {
		w       http.ResponseWriter
		status  int
		message string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				status: http.StatusBadRequest,
				w:      httptest.NewRecorder(),
			},
		},
		{
			name: "Test err",
			args: args{
				status: http.StatusBadRequest,
				w:      &errorWriter{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writeText(tt.args.w, tt.args.status, tt.args.message)
		})
	}
}

func TestHandler_listRoute(t *testing.T) {
	type fields struct {
		storage repository.Storager
		Router  *chi.Mux
		cfg     *sconfig.Config
	}
	tests := []struct {
		name    string
		request string
		fields  fields
	}{
		{
			name:    "Test 1",
			request: "/not-found",
		},
		{
			name:    "Test 2",
			request: "/ping",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			h := &Handler{
				storage: tt.fields.storage,
				Router:  chi.NewRouter(),
				cfg:     tt.fields.cfg,
			}

			h.listRoute()
			h.Router.ServeHTTP(w, r)
		})
	}
}
