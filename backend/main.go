// main.go
package main

import (
	"log"
	"net/http"
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
	surveyMux.HandleFunc("POST /api/surveys", handlers.CreateSurveyHandler) // Создание опроса
	surveyMux.HandleFunc("GET /api/surveys/analytics", func(w http.ResponseWriter, r *http.Request) {
		handlers.AnalyticsHandler(w, r)
	})

	surveyMux.HandleFunc("POST /api/surveys/submit", func(w http.ResponseWriter, r *http.Request) {
		handlers.SubmitSurveyHandler(w, r)
	})

	surveyMux.HandleFunc("GET /api/surveys/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.GetSurveyHandler(w, r)
	})
	surveyMux.HandleFunc("PUT /api/surveys/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.EditSurveyHandler(w, r)
	})

	// Применение middleware к защищённым маршрутам
	protectedSurveyHandler := handlers.AuthMiddleware(surveyMux)
	http.Handle("/api/surveys/", protectedSurveyHandler)

	// Запуск сервера
	log.Println("Сервер запущен на :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
