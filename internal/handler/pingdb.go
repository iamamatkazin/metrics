package handler

import (
	"net/http"
)

// pingDB - обрабатывает входной запрос на получение пинга от базы данных.
func (h *Handler) pingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.Ping(r.Context()); err != nil {
		writeText(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeText(w, http.StatusOK, http.StatusText(http.StatusOK))
}
