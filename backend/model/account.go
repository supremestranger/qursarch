package model

import (
	"backend/accounts"
	"backend/auth"
	"backend/utils"
	"encoding/json"
	"net/http"
)

type SignUpRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func RegisterAccountModels() {
	utils.RegisterOnGet("/accounts", onSignIn)
	utils.RegisterOnPost("/accounts", onSignUp)
}

func onSignIn(rw http.ResponseWriter, req *http.Request) {
	// todo проверить что пароль правильный

	username := ""
	token, err := auth.CreateToken(username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	cookie := http.Cookie{
		Name:  "token",
		Value: token,
	}
	http.SetCookie(rw, &cookie)
}

func onSignUp(rw http.ResponseWriter, req *http.Request) {
	var signUpReq SignUpRequest
	err := json.NewDecoder(req.Body).Decode(&signUpReq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	if len(signUpReq.Login) == 0 {
		http.Error(rw, "Too short username", http.StatusBadRequest)
	}

	if len(signUpReq.Password) == 0 {
		http.Error(rw, "Too short password", http.StatusBadRequest)
	} // todo проверить что ник не занят

	accounts.CreateAccount(accounts.AccountDesc{Login: signUpReq.Login, Password: signUpReq.Password})

	// todo записать в бд новый аккаунт
}
