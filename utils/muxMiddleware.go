package utils

import (
	"net/http"
	"strings"
)

func SubstituteXHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/_health" {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-Forwarded-Host") != "" {
			r.URL.Host = r.Header.Get("X-Forwarded-Host")
			forwardedUri := strings.Split(r.Header.Get("X-Forwarded-Uri"), "?")
			r.URL.Path = forwardedUri[0]
			if len(forwardedUri) > 1 {
				r.URL.RawQuery = forwardedUri[1]
			}
			r.Form = r.URL.Query()
		}
		next.ServeHTTP(w, r)
	})
}
