package middleware

import (
	"log"
	"net/http"
	"time"
)

type wrappedResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *wrappedResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logging(next http.Handler) http.Handler {
	log.Println("Logging is now enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		wrapped := &wrappedResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(wrapped, r)

		log.Printf("%s %s %s %d %s", r.Method, r.RequestURI, r.Proto, wrapped.statusCode, time.Since(start))
	})
}
