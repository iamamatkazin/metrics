package handler

import (
	"fmt"
	"net/http"
)

func ExampleHandler_listMetrics() {
	// получаем список метрик
	list := h.storage.ListMetrics()

	var li string
	for range list {
		// формируем в переменную li тело из списока метрик
	}

	// формируем html страницу
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

	// отсылаем ответ
	writeHTML(w, http.StatusOK, html)
}
