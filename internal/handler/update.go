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

	if err := h.storage.UpdateMetric(r.Context(), &metric); err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

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

	if err := h.storage.UpdateMetric(r.Context(), &metric); err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

func (h *Handler) updatesMetricJSON(w http.ResponseWriter, r *http.Request) {
	var metrics []model.Metric
	if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.UpdateMetrics(r.Context(), metrics); err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}
