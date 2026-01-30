package handler

import (
	"io"
	"net/http"
)

// gzipWriter - обертка над стандартным http.ResponseWriter.
type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

// Write - реализует метод интерфейса http.ResponseWriter.
func (w *gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
