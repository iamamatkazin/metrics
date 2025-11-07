package handler

import (
	"encoding/json"
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

	writeText(w, code, message)
}

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

	if value := h.storage.GetMetric(metric.ID); value != nil {
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

		writeJSON(w, http.StatusOK, body)
	} else {
		writeText(w, http.StatusNotFound, http.StatusText(http.StatusNotFound))
		// writeJSON(w, http.StatusNotFound, []byte("{\"error\": \"Not Found\"}"))
	}
}
