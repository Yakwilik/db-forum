package middleware

import "net/http"

func ContentTypeMiddleware(contentType string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		next.ServeHTTP(w, r)
	})

}
