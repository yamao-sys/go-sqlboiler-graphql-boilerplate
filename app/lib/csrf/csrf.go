package csrf

import (
	"net/http"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: Originが一致していない場合は403を返す
		origin := r.Header.Get("Origin")
		if origin != "http://localhost:3000" && origin != "http://localhost:8080" {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
