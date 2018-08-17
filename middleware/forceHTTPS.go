package middleware

import (
	"net/http"
)

// ForceHTTPS is middleware which redirects the user to https if the x-forward-proto header is set to http
func ForceHTTPS(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		proto := r.Header.Get("x-forwarded-proto")
		if proto == "http" {
			http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusPermanentRedirect)
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
