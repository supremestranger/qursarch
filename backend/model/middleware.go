package model

import (
	"backend/auth"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) (bool, string) {
	token, err := r.Cookie("token")
	if err != nil {
		tokenHeader, ok := r.Header["Authorization"]
		if !ok {
			return false, ""
		}
		return auth.VerifyToken(tokenHeader[0])
	}
	return auth.VerifyToken(token.Value)
}
