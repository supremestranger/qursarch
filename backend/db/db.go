// db/db.go
package db

import (
    "database/sql"
    "fmt"
    "log"
    "os"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    dbHost := os.Getenv("DB_HOST")
    dbPort := os.Getenv("DB_PORT")
    dbUser := os.Getenv("DB_USER")
    dbPassword := os.Getenv("DB_PASSWORD")
    dbName := os.Getenv("DB_NAME")

    psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
        "password=%s dbname=%s sslmode=disable",
        dbHost, dbPort, dbUser, dbPassword, dbName)

    DB, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Не удалось пинговать базу данных: %v", err)
    }

    log.Println("Успешно подключились к базе данных")
}
