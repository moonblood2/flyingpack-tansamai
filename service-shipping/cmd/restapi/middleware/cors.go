package middleware

import (
	"net/http"
)

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Headers", "Authorization, authentication, Content-Type")
		header.Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE")
		//Handle preflight request.
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
		} else { //Handle actual request.
			next.ServeHTTP(w, r)
		}
	})
}
