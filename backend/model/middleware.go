package model

import (
	"backend/auth"
	"net/http"
)

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := auth.VerifyToken(token.Value)

		if ok {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "проверьте логин или пароль", http.StatusForbidden)
		}
	})
}
