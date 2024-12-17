package model

import (
	"backend/auth"
	"log"
	"net/http"
)

func CheckAuth(w http.ResponseWriter, r *http.Request) (bool, string) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Println("cookie!")
		tokenHeader, ok := r.Header["Authorization"]
		if !ok {
			return false, ""
		}
		return auth.VerifyToken(tokenHeader[0])
	}
	log.Println("fwefwef")
	return auth.VerifyToken(token.Value)
}
