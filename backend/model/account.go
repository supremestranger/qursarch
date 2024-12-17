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
	utils.RegisterOnPost(ACCOUNT_ROOT+"/check_auth", onCheckAuth)
	utils.RegisterOnGet(ACCOUNT_ROOT, onSignIn)
	utils.RegisterOnPost(ACCOUNT_ROOT, onSignUp)
}

func onCheckAuth(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw, "http://localhost:3000")
	rw.Header().Add("Access-Control-Allow-Credentials", "true")
	ok, err := CheckAuth(rw, req)
	if !ok {
		log.Println(err)
		http.Error(rw, "вы не авторизованы", http.StatusBadRequest)
		return
	}
}

func onSignIn(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw, "http://localhost:3000")
	// todo проверить что пароль правильный

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
	// ok := accounts.Login(signUpReq.Login, signUpReq.Password)
	token, err := auth.CreateToken(signUpReq.Login)
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
	utils.EnableCors(rw, "http://localhost:3000")
	rw.Header().Add("Access-Control-Allow-Credentials", "true")
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
		Path:  "/", // Cookie available for the entire site
	}
	http.SetCookie(rw, &cookie)
}
