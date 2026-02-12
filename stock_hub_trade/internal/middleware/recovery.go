package middleware

import (
	"log"
	"net/http"
	"runtime/debug"
)

// Recovery recovers from panics and returns 500 error
func (m *Manager) Recovery(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("❌ PANIC: %v\n%s", err, debug.Stack())

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"Internal server error"}`))
			}
		}()

		next(w, r)
	}
}
