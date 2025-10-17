package handler

import "net/http"

type responseData struct {
	status int
	size   int
}

type responseWriter struct {
	http.ResponseWriter
	data *responseData
}

func (r *responseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.data.size += size

	return size, err
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.data.status = statusCode
}
