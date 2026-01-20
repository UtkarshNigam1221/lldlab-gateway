package middleware

import (
	"net/http"
	"strings"
)

// AuthMiddleware handles authentication for incoming requests.
// Currently passes through all requests but can be extended for JWT validation.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if r.URL.Path == "/health" {
			next.ServeHTTP(w, r)
			return
		}

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(authHeader, "Bearer ") {
			next.ServeHTTP(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}
