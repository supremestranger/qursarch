// db/db.go
package db

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
    var err error
    connStr := "user=yourusername password=yourpassword dbname=surveydb sslmode=disable"
    DB, err = sql.Open("postgres", connStr)
    if err != nil {
        log.Fatalf("Не удалось подключиться к базе данных: %v", err)
    }

    err = DB.Ping()
    if err != nil {
        log.Fatalf("Не удалось пинговать базу данных: %v", err)
    }

    fmt.Println("Успешно подключились к базе данных")
}
