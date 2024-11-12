package accounts

import (
	"backend/db"
	"errors"
	"log"
)

type AccountDesc struct {
	Login    string
	Password string
}

func CreateAccount(accountDesc AccountDesc) error {
	rows, err := db.DB.Query("SELECT * FROM ACCOUNTS WHERE ACCOUNTS.Username = $1", accountDesc.Login)
	for rows.Next() {
		return errors.New("ник занят")
	}
	result, err := db.DB.Exec("INSERT INTO ACCOUNTS (Username, Password) values ($1, $2)", accountDesc.Login, accountDesc.Password)
	if err != nil {
		return err
	}
	log.Println(result)
	return nil
}
