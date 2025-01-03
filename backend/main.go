// main.go
package main

import (
	"log"
	"net/http"
	"strconv"
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

	// Создание отдельного ServeMux для маршрутов опросов
	surveyMux := http.NewServeMux()

	// Обработчик для создания опроса (POST /api/surveys)
	surveyMux.HandleFunc("/api/surveys", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handlers.CreateSurveyHandler(w, r)
		case http.MethodGet:
			handlers.ListSurveysHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Обработчик для всех маршрутов, начинающихся с /api/surveys/
	surveyMux.HandleFunc("/api/surveys/", func(w http.ResponseWriter, r *http.Request) {
		// Извлечение части URL после /api/surveys/
		path := strings.TrimPrefix(r.URL.Path, "/api/surveys/")
		if path == "" {
			http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
			return
		}

		// Разделение пути на части
		parts := strings.SplitN(path, "/", 2)
		surveyIDStr := parts[0]
		surveyID, err := strconv.Atoi(surveyIDStr)
		_ = surveyID
		if err != nil {
			http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
			return
		}

		if len(parts) == 1 {
			// Маршруты вида /api/surveys/{id}
			if r.Method == http.MethodGet {
				handlers.GetSurveyHandler(w, r)
				return
			} else if r.Method == http.MethodPut {
				handlers.EditSurveyHandler(w, r)
				return
			}
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		// Маршруты вида /api/surveys/{id}/analytics или /api/surveys/{id}/submit
		subPath := parts[1]
		switch subPath {
		case "analytics":
			if r.Method == http.MethodGet {
				handlers.AnalyticsHandler(w, r)
				return
			}
		case "submit":
			if r.Method == http.MethodPost {
				handlers.SubmitSurveyHandler(w, r)
				return
			}
		}

		// Если маршрут не соответствует ни одному из вышеуказанных
		http.Error(w, "Неизвестный маршрут", http.StatusNotFound)
	})

	// Применение AuthMiddleware к маршрутам опросов
	protectedSurveyHandler := handlers.AuthMiddleware(surveyMux)
	deleteSurveyHandler := handlers.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// todo delete survey
		path := strings.TrimPrefix(r.URL.Path, "/api/surveys_delete/")
		if path == "" {
			http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
			return
		}

		// Разделение пути на части
		parts := strings.SplitN(path, "/", 2)
		surveyIDStr := parts[0]
		surveyID, err := strconv.Atoi(surveyIDStr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res, err := db.DB.Exec("DELETE FROM surveys WHERE surveys.SurveyID = ($1)", surveyID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, _ = res.RowsAffected()
	}))

	// Регистрация защищённых маршрутов
	http.Handle("/api/surveys", protectedSurveyHandler)  // Для точного соответствия /api/surveys
	http.Handle("/api/surveys/", protectedSurveyHandler) // Для всех маршрутов, начинающихся с /api/surveys/
	http.HandleFunc("/api/public_surveys/", handlers.GetSurveyHandler)
	http.HandleFunc("/api/public_surveys_submit/", handlers.SubmitSurveyHandler)
	http.Handle("/api/surveys_delete/", deleteSurveyHandler)
	// Запуск сервера
	log.Println("Сервер запущен на :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
