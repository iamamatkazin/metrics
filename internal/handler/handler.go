package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/iamamatkazin/metrics.git/internal/repository"
)

type Handler struct {
	storage repository.Storager
	Router  *chi.Mux
}

func New() *Handler {
	h := &Handler{
		storage: repository.New(),
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h
}

func (h *Handler) listRoute() {
	h.Router.Use(middlewareLog)
	h.Router.Get("/", h.listMetrics)
	h.Router.Get("/value/{type}/{id}", h.getMetric)
	h.Router.Post("/update/{type}/{id}/{val}", h.updateMetric)

	h.Router.With(middleware.AllowContentType("application/json")).Post("/value/", h.getMetricJSON)
	h.Router.With(middleware.AllowContentType("application/json")).Post("/update/", h.updateMetricJSON)

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)

		if _, err := w.Write([]byte(http.StatusText(http.StatusNotFound))); err != nil {
			slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
		}
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)

		if _, err := w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed))); err != nil {
			slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
		}
	})
}

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			data:           &responseData{},
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		slog.Info("Запрос",
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.Duration("duration", duration),
			slog.Int("status", rw.data.status),
			slog.Int("size", rw.data.size),
		)
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
