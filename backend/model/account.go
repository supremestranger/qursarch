package model

import (
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
	// todo проверить что ник не занят

	// todo записать в бд новый аккаунт
}
