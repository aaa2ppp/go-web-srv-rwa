package http_debug

import (
	"log"
	"net/http"
)

type tDebugWriter struct {
	http.ResponseWriter
	status int
}

func newDebugWriter(w http.ResponseWriter) *tDebugWriter {
	return &tDebugWriter{
		ResponseWriter: w,
		status:         200,
	}
}

func (w *tDebugWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func Handle(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		dw := newDebugWriter(w)
		next.ServeHTTP(dw, r)
		log.Printf("[%s] %s - %d %s", r.Method, r.URL.Path, dw.status, http.StatusText(dw.status))
	}
}
