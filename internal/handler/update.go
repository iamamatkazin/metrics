package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (h *Handler) updateMetric(w http.ResponseWriter, r *http.Request) {
	metric := model.Metric{
		ID:    chi.URLParam(r, "id"),
		MType: chi.URLParam(r, "type"),
	}

	if err := metric.Validate(); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := metric.Normalize(chi.URLParam(r, "val")); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	h.storage.UpdateMetric(metric)
	writeText(w, http.StatusOK, http.StatusText(http.StatusOK))
}

func (h *Handler) updateMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metric model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metric); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := metric.ValidateJSON(); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	h.storage.UpdateMetric(metric)
	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}
