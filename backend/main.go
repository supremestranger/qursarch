// main.go
package main

import (
    "log"
    "net/http"
    "strings"

    "survey-platform-server/db"
    "survey-platform-server/handlers"
    "survey-platform-server/middleware"
)

func main() {
    // Инициализация базы данных
    db.InitDB()

    // Регистрация обработчиков маршрутов
    http.HandleFunc("/register", handlers.RegisterHandler)
    http.HandleFunc("/login", handlers.LoginHandler)
    http.HandleFunc("/logout", handlers.LogoutHandler)
    http.HandleFunc("/check_auth", handlers.CheckAuthHandler)

    // Защищённые маршруты
    surveyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/survey/create") {
            handlers.CreateSurveyHandler(w, r)
            return
        }
        if strings.HasPrefix(r.URL.Path, "/survey/edit/") {
            handlers.EditSurveyHandler(w, r)
            return
        }
        if strings.HasSuffix(r.URL.Path, "/submit") {
            handlers.SubmitSurveyHandler(w, r)
            return
        }
        if r.Method == http.MethodGet {
            handlers.GetSurveyHandler(w, r)
            return
        }
        http.Error(w, "Неизвестный маршрут", http.StatusNotFound)
    })

    analyticsHandler := http.HandlerFunc(handlers.AnalysisHandler)

    // Применение middleware к защищённым маршрутам
    http.Handle("/survey/", middleware.AuthMiddleware(surveyHandler))
    http.Handle("/analytics/", middleware.AuthMiddleware(analyticsHandler))

    // Запуск сервера
    log.Println("Сервер запущен на :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
