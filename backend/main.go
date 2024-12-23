// main.go
package main

import (
    "log"
    "net/http"
    "strings"
    "survey-platform-server/db"
    "survey-platform-server/handlers"
)

func main() {
    // Инициализация базы данных
    db.InitDB()

    // Маршруты аутентификации
    http.HandleFunc("/api/register", handlers.RegisterHandler)
    http.HandleFunc("/api/login", handlers.LoginHandler)
    http.HandleFunc("/api/logout", handlers.LogoutHandler)
    http.HandleFunc("/api/check_auth", handlers.CheckAuthHandler)

    // Защищённые маршруты для управления опросами и аналитики
    // Используем AuthMiddleware для защиты этих маршрутов
    surveyMux := http.NewServeMux()
    surveyMux.HandleFunc("/api/surveys", handlers.CreateSurveyHandler) // Создание опроса
    surveyMux.HandleFunc("/api/surveys/", func(w http.ResponseWriter, r *http.Request) {
        if strings.HasSuffix(r.URL.Path, "/analytics") {
            handlers.AnalyticsHandler(w, r)
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
        if r.Method == http.MethodPut {
            handlers.EditSurveyHandler(w, r)
            return
        }
        http.Error(w, "Неизвестный маршрут", http.StatusNotFound)
    })

    // Применение middleware к защищённым маршрутам
    protectedSurveyHandler := handlers.AuthMiddleware(surveyMux)
    http.Handle("/api/surveys/", protectedSurveyHandler)

    // Запуск сервера
    log.Println("Сервер запущен на :8081")
    log.Fatal(http.ListenAndServe(":8081", nil))
}
