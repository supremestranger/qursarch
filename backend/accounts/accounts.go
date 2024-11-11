package accounts

import (
	"backend/db"
	"log"
)

type AccountDesc struct {
	Login    string
	Password string
}

func CreateAccount(accountDesc AccountDesc) {
	result, err := db.DB.Exec("INSERT INTO ACCOUNTS (Username, Password) values ($1, $2)", accountDesc.Login, accountDesc.Password)
	if err != nil {
		log.Println(err)
	}
	log.Println(result)
}
