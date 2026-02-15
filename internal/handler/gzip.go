package handler

import (
	"io"
	"net/http"
)

// gzipWriter оборачивает http.ResponseWriter для сжатия ответов gzip.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write реализует метод интерфейса http.ResponseWriter.
func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
