package model

import (
	"backend/accounts"
	"backend/auth"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

const ACCOUNT_ROOT = "/accounts"

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterAccountModels() {
	utils.RegisterOnGet(ACCOUNT_ROOT, onSignIn)
	utils.RegisterOnPost(ACCOUNT_ROOT, onSignUp)
}

func onSignIn(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw)
	// todo проверить что пароль правильный

	username := ""
	token, err := auth.CreateToken(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(rw, &cookie)
}

func onSignUp(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw)
	var signUpReq SignUpRequest
	err := json.NewDecoder(req.Body).Decode(&signUpReq)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	if len(signUpReq.Login) == 0 {
		log.Println(err)
		http.Error(rw, "Too short username", http.StatusBadRequest)
		return
	}

	if len(signUpReq.Password) == 0 {
		log.Println(err)
		http.Error(rw, "Too short password", http.StatusBadRequest)
		return
	}

	err = accounts.CreateAccount(accounts.AccountDesc{Login: signUpReq.Login, Password: signUpReq.Password})
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := auth.CreateToken(signUpReq.Login)
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(rw, &cookie)
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
}
