package middleware

import (
	"burakozkan138/questionanswerapi/internal/models"
	"log"
	"net/http"
)

func ErrorHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				response := models.NewResponse(false, "Internal Server Error", http.StatusInternalServerError, nil, nil)
				response.Write(w)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
