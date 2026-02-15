// Package handler предоставляет HTTP обработчики для сервера метрик.
// Реализует Chi роутер с эндпоинтами для получения, хранения и извлечения метрик.
// Пакет поддерживает JSON и URL-encoded форматы запросов, а также gzip сжатие ответов.
package handler

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/internal/observer"
	"github.com/iamamatkazin/metrics.git/internal/pool"
	"github.com/iamamatkazin/metrics.git/internal/repository"
	sconfig "github.com/iamamatkazin/metrics.git/pkg/config/server"
)

// Handler - главная структура приложения.
// generate:reset
type Handler struct {
	storage     repository.Storager
	Router      *chi.Mux
	cfg         *sconfig.Config
	audit       *observer.Event
	poolMessage *pool.Pool[*model.Message]
}

// New - конструктор для Handler.
func New(ctx context.Context, cfg *sconfig.Config) (*Handler, error) {
	storage, err := repository.New(ctx, cfg)
	if err != nil {
		return nil, err
	}

	audit, err := observer.New(cfg)
	if err != nil {
		return nil, err
	}

	h := &Handler{
		storage: storage,
		cfg:     cfg,
		audit:   audit,
		poolMessage: pool.New(func() *model.Message {
			return &model.Message{}
		}),
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h, nil
}

// listRoute - формирует список роутов.
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

	h.Router.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(http.StatusText(http.StatusNotFound)))
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	})
}

// writeText - формирует ответ сервера в виде текста.
func writeText(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(message)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

// writeJSON - формирует ответ сервера в виде json структуры.
func writeJSON(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

// writeHTML - формирует ответ сервера в виде html структуры.
func writeHTML(w http.ResponseWriter, status int, html string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}

// Shutdown - коректное завершение сервиса.
func (h *Handler) Shutdown() {
	h.storage.Shutdown()
}
