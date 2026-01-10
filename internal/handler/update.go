package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/common"
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

	if err := checkSign(r, h.cfg.Key, metrics); err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.storage.UpdateMetrics(r.Context(), metrics); err != nil {
		writeText(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

func checkSign(r *http.Request, key string, body any) error {
	if key == "" {
		return nil
	}

	data, err := json.Marshal(body)
	if err != nil {
		return err
	}

	calcValue := common.CalcSign([]byte(key), data)
	headerValue := r.Header.Get("HashSHA256")

	if !reflect.DeepEqual(calcValue, headerValue) {
		return fmt.Errorf("подпись не верна")
	}

	return nil
}
