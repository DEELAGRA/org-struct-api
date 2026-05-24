package middleware

import (
	"log"
	"net/http"
	"time"
)

type responseWrite struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWrite) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapped := &responseWrite{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)
		log.Printf(
			"%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start).Round(time.Microsecond),
		)
	})
}
