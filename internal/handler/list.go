package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/iamamatkazin/metrics.git/internal/model"
)

func (h *Handler) listMetrics(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte(html)); err != nil {
		slog.Error("Ошибка отправки ответа:", slog.Any("error", err))
	}
}
