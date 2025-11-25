package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamamatkazin/metrics.git/internal/repository"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

type Handler struct {
	storage repository.Storager
	Router  *chi.Mux
	cfg     *sconfig.Config
}

func New(ctx context.Context, cfg *sconfig.Config) (*Handler, error) {
	storage, err := repository.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		storage: storage,
		cfg:     cfg,
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h, nil
}

func (h *Handler) listRoute() {
	h.Router.Use(middlewareLog)
	h.Router.Use(middlewareGzip)
	h.Router.Get("/ping", h.pingDB)
	h.Router.Get("/", h.listMetrics)
	h.Router.Get("/value/{type}/{id}", h.getMetric)
	h.Router.Post("/update/{type}/{id}/{val}", h.updateMetric)

	h.Router.With(middleware.AllowContentType("application/json")).Post("/value/", h.getMetricJSON)
	h.Router.With(middleware.AllowContentType("application/json")).Post("/update/", h.updateMetricJSON)
	h.Router.With(middleware.AllowContentType("application/json")).Post("/updates/", h.updatesMetricJSON)

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})
}

func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(message)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func writeHTML(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

func (h *Handler) Shutdown() {
	h.storage.Shutdown()
}
