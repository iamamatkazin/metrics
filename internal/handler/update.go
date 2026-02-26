package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/common"
	"github.com/iamamatkazin/metrics.git/internal/model"
	"github.com/iamamatkazin/metrics.git/pkg/crypto"
)

// updateMetric сохраняет метрику из строки запроса.
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

	if message := h.getMessage([]model.Metric{metric}, r.RemoteAddr); message != nil {
		h.audit.Send(*message)
		h.poolMessage.Put(message)
	}

	writeText(w, http.StatusOK, http.StatusText(http.StatusOK))
}

// updateMetricJSON сохраняет метрику из тела запроса в формате JSON.
func (h *Handler) updateMetricJSON(w http.ResponseWriter, r *http.Request) {
	decodeBody, err := getDecodeBody(r, h.cfg.CryptoKey)
	if err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	var metric model.Metric
	if err := json.Unmarshal(decodeBody, &metric); err != nil {
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

	if message := h.getMessage([]model.Metric{metric}, r.RemoteAddr); message != nil {
		h.audit.Send(*message)
		h.poolMessage.Put(message)
	}

	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

// updatesMetricJSON сохраняет массив метрик из тела запроса в формате JSON.
func (h *Handler) updatesMetricJSON(w http.ResponseWriter, r *http.Request) {
	decodeBody, err := getDecodeBody(r, h.cfg.CryptoKey)
	if err != nil {
		writeText(w, http.StatusBadRequest, err.Error())
		return
	}

	var metrics []model.Metric
	if err := json.Unmarshal(decodeBody, &metrics); err != nil {
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

	if message := h.getMessage(metrics, r.RemoteAddr); message != nil {
		h.audit.Send(*message)
		h.poolMessage.Put(message)
	}

	writeJSON(w, http.StatusOK, []byte("{\"status\": \"OK\"}"))
}

// checkSign проверяет HMAC-SHA256 подпись запроса.
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
		return fmt.Errorf("подпись неверна")
	}

	return nil
}

// getMessage создает сообщение для аудита.
func (h *Handler) getMessage(list []model.Metric, ip string) *model.Message {
	if len(list) == 0 {
		return nil
	}

	metrics := make([]string, 0, len(list))
	for _, item := range list {
		metrics = append(metrics, item.ID)
	}

	message := h.poolMessage.Get()
	message.Date = time.Now().Unix()
	message.Metrics = metrics
	message.IP = ip

	return message
}

func getDecodeBody(r *http.Request, cryptoKey string) ([]byte, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	decodeBody, err := crypto.Decrypt(cryptoKey, body)
	if err != nil {
		return nil, err
	}

	return decodeBody, nil
}
