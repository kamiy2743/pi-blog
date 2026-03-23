package middleware

import (
	"net/http"
	"strings"
)

func NormalizePath() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			path := r.URL.Path
			if path != "/" && strings.HasSuffix(path, "/") {
				normalizedPath := strings.TrimRight(path, "/")
				if r.URL.RawQuery != "" {
					normalizedPath += "?" + r.URL.RawQuery
				}

				http.Redirect(w, r, normalizedPath, http.StatusPermanentRedirect)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
