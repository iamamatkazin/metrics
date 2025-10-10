package handler

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/iamamatkazin/metrics.git/internal/repository"
)

type Handler struct {
	storage repository.Storager
	Router  *chi.Mux
}

func New() *Handler {
	h := &Handler{
		storage: repository.New(),
	}

	h.Router = chi.NewRouter()
	h.listRoute()

	return h
}

func (h *Handler) listRoute() {
	h.Router.Get("/", h.listMetrics)
	h.Router.Post("/update/{type}/{id}/{val}", h.updateMetric)
	h.Router.Get("/value/{type}/{id}", h.getMetric)

	h.Router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, http.StatusText(http.StatusNotFound))
	})

	h.Router.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, http.StatusText(http.StatusMethodNotAllowed))
	})
}
