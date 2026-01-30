package handler

import (
	"fmt"
	"net/http"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

// listMetrics обрабатывает входной запрос и возвращает HTML страницу со списком метрик.
func (h *Handler) listMetrics(w http.ResponseWriter, _ *http.Request) {
	var (
		li  string
		val any
	)

	list := h.storage.ListMetrics()

	for i := range list {
		if list[i].MType == model.Counter {
			val = *list[i].Delta
		} else {
			val = *list[i].Value
		}
		li += fmt.Sprintf("<li>%s: %v</li>", list[i].ID, val)
	}

	html := fmt.Sprintf(`
	<html>
		<head>
		<title></title>
		</head>
		<body>
			<ul>
			%s
			</ul>
		</body>
	</html>`, li)

	writeHTML(w, http.StatusOK, html)
}
