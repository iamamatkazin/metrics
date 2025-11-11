package handler

import (
	"net/http"
)

func (h *Handler) pingDB(w http.ResponseWriter, r *http.Request) {
	if err := h.storage.PingDB(r.Context()); err != nil {
		writeText(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	writeText(w, http.StatusOK, http.StatusText(http.StatusOK))
}
