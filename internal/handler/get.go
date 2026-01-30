package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/common"
	"github.com/iamamatkazin/metrics.git/internal/model"
)

// getMetric - обрабатывает входной запрос на получение значения метрики по ее идентификатору.
func (h *Handler) getMetric(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	message := http.StatusText(http.StatusNotFound)

	value, err := h.storage.GetMetric(r.Context(), chi.URLParam(r, "id"))
	if err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

	if value != nil {
		code = http.StatusOK

		if value.MType == model.Gauge {
			message = strconv.FormatFloat(*value.Value, 'f', -1, 64)
		} else {
			message = strconv.Itoa(*value.Delta)
		}
	}

	writeText(w, code, message)
}

// getMetricJSON - обрабатывает входной запрос на получение метрики по ее идентификатору.
func (h *Handler) getMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := metric.Validate(); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	value, err := h.storage.GetMetric(r.Context(), metric.ID)
	if err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

	if value != nil {
		if metric.MType == model.Gauge {
			metric.Value = value.Value
		} else {
			metric.Delta = value.Delta
		}

		body, err := json.Marshal(metric)
		if err != nil {
			writeText(w, http.StatusInternalServerError, err.Error())
			return
		}

		if h.cfg.Key != "" {
			w.Header().Set("HashSHA256", string(common.CalcSign([]byte(h.cfg.Key), body)))
		}

		writeJSON(w, http.StatusOK, body)
	} else {
		writeText(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
	}
}
