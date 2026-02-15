package middleware

import (
	"net/http"
	"strings"
)

// CORSConfig holds the CORS policy settings.
type CORSConfig struct {
	AllowedOrigins []string // e.g. ["https://myapp.com", "http://localhost:3000"]
	AllowedMethods []string // e.g. ["GET", "POST", "PUT", "PATCH", "DELETE"]
	AllowedHeaders []string // e.g. ["Content-Type", "X-API-Key"]
	MaxAge         string   // preflight cache duration in seconds, e.g. "86400"
}

// DefaultCORSConfig returns a restrictive but practical default config.
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: []string{"http://localhost:3000"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-API-Key"},
		MaxAge:         "86400", // 24 hours
	}
}

// CORSMiddleware handles cross-origin resource sharing.
// It checks the request Origin against the AllowedOrigins whitelist.
// Preflight OPTIONS requests are responded to immediately (no auth needed).
func CORSMiddleware(cfg CORSConfig, next http.Handler) http.Handler {
	allowedOrigins := make(map[string]bool, len(cfg.AllowedOrigins))
	for _, o := range cfg.AllowedOrigins {
		allowedOrigins[o] = true
	}

	methods := strings.Join(cfg.AllowedMethods, ", ")
	headers := strings.Join(cfg.AllowedHeaders, ", ")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Only set CORS headers when an Origin is present
		if origin != "" {
			if allowedOrigins[origin] || allowedOrigins["*"] {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin") // tell caches the response varies by Origin
			} else {
				// Origin not allowed — reject with 403
				http.Error(w, `{"error":"origin not allowed"}`, http.StatusForbidden)
				return
			}

			w.Header().Set("Access-Control-Allow-Methods", methods)
			w.Header().Set("Access-Control-Allow-Headers", headers)
			w.Header().Set("Access-Control-Max-Age", cfg.MaxAge)
		}

		// Handle preflight — respond immediately, skip auth & handler
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
