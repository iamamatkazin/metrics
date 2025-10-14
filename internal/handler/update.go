package handler

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	metric := model.Metric{
		ID:    chi.URLParam(r, "id"),
		MType: chi.URLParam(r, "type"),
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	if err := metric.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	if err := metric.Normalize(chi.URLParam(r, "val")); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())
		return
	}

	h.storage.UpdateMetric(metric)
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(http.StatusText(http.StatusOK))); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}
