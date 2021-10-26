package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
)

func HTTPSetMaxFileSize(maxSizeMb int64) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.MaxBytesReader(w, r.Body, maxSizeMb * 1024 * 1024)
			next.ServeHTTP(w, r)
		})
	}
}