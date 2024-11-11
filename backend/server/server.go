package server

import (
	"backend/db"
	"backend/model"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	host     = "database"
	port     = 5432
	user     = "QDBSU"
	password = "TheSacredKailash"
	dbname   = "QDB"
)

func Log() {
	log.Println("server working")
}

func InitRouting() {
	model.RegisterAccountModels()
	model.RegisterSurveyModels()
}

func InitDb() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Println(err)
	}
	sqlErr := db.DB.Ping()
	if sqlErr != nil {
		fmt.Println(sqlErr)
	}
}
