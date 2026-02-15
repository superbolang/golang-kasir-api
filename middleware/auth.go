package middleware

import (
	"encoding/json"
	"net/http"
)

// APIKeyMiddleware validates the API key on every request.
// The key is accepted from two places (in order of priority):
//  1. HTTP Header  →  X-API-Key: <key>
//  2. Query param  →  ?api_key=<key>
func APIKeyMiddleware(validKey string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("X-API-Key")

		// Fallback: check query parameter
		if key == "" {
			key = r.URL.Query().Get("api_key")
		}

		if key == "" {
			respondUnauthorized(w, "missing API key")
			return
		}

		if !validateKey(key, validKey) {
			respondUnauthorized(w, "invalid API key")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// validateKey does a constant-time comparison to prevent timing attacks.
func validateKey(provided, valid string) bool {
	if len(provided) != len(valid) {
		return false
	}
	// XOR every byte — result is 0 only when all bytes match
	var diff byte
	for i := 0; i < len(provided); i++ {
		diff |= provided[i] ^ valid[i]
	}
	return diff == 0
}

func respondUnauthorized(w http.ResponseWriter, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}

// LoggingMiddleware logs every incoming request (method, path, remote addr).
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Mask the api_key value in logged URLs
		logURL := r.URL.Path
		if r.URL.RawQuery != "" {
			q := r.URL.Query()
			if q.Get("api_key") != "" {
				q.Set("api_key", "***")
			}
			logURL += "?" + q.Encode()
		}
		println("[" + r.Method + "] " + logURL + " from " + r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

// Chain applies multiple middlewares to a handler (outermost first).
// Usage: Chain(handler, Logging, Auth) → Logging → Auth → handler
func Chain(h http.Handler, middlewares ...func(http.Handler) http.Handler) http.Handler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		h = middlewares[i](h)
	}
	return h
}
