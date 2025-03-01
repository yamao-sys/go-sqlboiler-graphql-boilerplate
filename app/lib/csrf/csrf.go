package csrf

import (
	"net/http"
	"os"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE: Originが一致していない場合は403を返す
		origin := r.Header.Get("Origin")
		if origin != os.Getenv("CLIENT_ORIGIN") && origin != os.Getenv("ORIGIN") {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
