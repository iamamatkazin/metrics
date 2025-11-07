package handler

import (
	"compress/gzip"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func middlewareLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			data:           &responseData{},
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start)
		slog.Info("Запрос",
			slog.String("uri", r.RequestURI),
			slog.String("method", r.Method),
			slog.Duration("duration", duration),
			slog.Int("status", rw.data.status),
			slog.Int("size", rw.data.size),
		)
	})
}

func middlewareGzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			writeText(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer gz.Close()

		gzw := &gzipWriter{ResponseWriter: w, Writer: gz}

		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzw, r)
	})
}
