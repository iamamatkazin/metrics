package handler

import "net/http"

// responseData хранит статус и размер HTTP ответа.
type responseData struct {
	status int
	size   int
}

// responseWriter оборачивает http.ResponseWriter для сбора статистики ответа.
type responseWriter struct {
	http.ResponseWriter
	data *responseData
}

// Write реализует метод интерфейса http.ResponseWriter.
func (r *responseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.data.size += size

	return size, err
}

// WriteHeader реализует метод интерфейса http.ResponseWriter.
func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.data.status = statusCode
}
