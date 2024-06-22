package middleware

import (
	"log"
	"net/http"
)

type Middleware func(http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			x := middlewares[i]
			next = x(next)
		}
		return next
	}
}

func AllowCors(next http.Handler) http.Handler {
	log.Println("Cors Enabled")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Origin")

		next.ServeHTTP(w, r)
	})
}
