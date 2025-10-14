package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	message := http.StatusText(http.StatusNotFound)

	if value := h.storage.GetMetric(chi.URLParam(r, "id")); value != nil {
		code = http.StatusOK

		if value.MType == model.Gauge {
			message = strconv.FormatFloat(*value.Value, 'f', -1, 64)
		} else {
			message = strconv.Itoa(*value.Delta)
		}
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	if _, err := w.Write([]byte(message)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}
