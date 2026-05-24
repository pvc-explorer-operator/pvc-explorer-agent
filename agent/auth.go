package agent

import (
	"crypto/subtle"
	"net/http"
)

// WithAuth returns a handler that enforces HTTP Bearer token authentication.
// If token is empty the handler passes through — this is used for standalone
// mode where no authentication is required.
//
// All routes behind this middleware are protected when a token is configured:
// any request without a valid Authorization: Bearer <token> header receives a
// 401 response.
func WithAuth(next http.Handler, token string) http.Handler {
	if token == "" {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if len(auth) < 7 || auth[:7] != "Bearer " {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		given := auth[7:]

		// Constant-time compare prevents timing side-channels.
		if subtle.ConstantTimeCompare([]byte(given), []byte(token)) != 1 {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
