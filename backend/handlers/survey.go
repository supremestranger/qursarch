// handlers/survey.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"survey-platform-server/db"
	"survey-platform-server/models"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetSurveysHandler(w http.ResponseWriter, r *http.Request) {
	var res []models.Survey
	rows, err := db.DB.Query("SELECT * FROM surveys")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for rows.Next() {
		var survey models.Survey
		if err := rows.Scan(&survey.SurveyID, &survey.Title, &survey.Description, &survey.CreatedAt, &survey.CreatedBy); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		res = append(res, survey)
	}

	json.NewEncoder(w).Encode(res)
}

// CreateSurveyHandler создает новый опрос
func CreateSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получение AdminID из контекста
	adminID, ok := GetAdminID(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить AdminID", http.StatusInternalServerError)
		return
	}

	var survey models.Survey
	err := json.NewDecoder(r.Body).Decode(&survey)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if survey.Title == "" {
		http.Error(w, "Название опроса обязательно", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	var surveyID int
	err = tx.QueryRow(
		"INSERT INTO Surveys (Title, Description, CreatedBy) VALUES ($1, $2, $3) RETURNING SurveyID",
		survey.Title, survey.Description, adminID).Scan(&surveyID)
	if err != nil {
		http.Error(w, "Ошибка при создании опроса", http.StatusInternalServerError)
		return
	}

	for _, q := range survey.Questions {
		var questionID int
		err = tx.QueryRow(
			"INSERT INTO Questions (SurveyID, QuestionText, QuestionType) VALUES ($1, $2, $3) RETURNING QuestionID",
			surveyID, q.QuestionText, q.QuestionType).Scan(&questionID)
		if err != nil {
			http.Error(w, "Ошибка при добавлении вопроса", http.StatusInternalServerError)
			return
		}

		if q.QuestionType == "single_choice" || q.QuestionType == "multiple_choice" {
			for _, opt := range q.Options {
				_, err = tx.Exec(
					"INSERT INTO AnswerOptions (QuestionID, OptionText) VALUES ($1, $2)",
					questionID, opt.OptionText)
				if err != nil {
					http.Error(w, "Ошибка при добавлении варианта ответа", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Ошибка при сохранении опроса", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   "Опрос успешно создан",
		"survey_id": surveyID,
	})
}

func ListSurveysHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Выполнение запроса к базе данных для получения всех опросов
	rows, err := db.DB.Query("SELECT SurveyID, Title, Description, CreatedAt, CreatedBy FROM Surveys")
	if err != nil {
		http.Error(w, "Ошибка при получении опросов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var surveys []models.Survey

	for rows.Next() {
		var survey models.Survey
		var createdBy sql.NullInt32
		err := rows.Scan(&survey.SurveyID, &survey.Title, &survey.Description, &survey.CreatedAt, &createdBy)
		if err != nil {
			http.Error(w, "Ошибка при сканировании опросов", http.StatusInternalServerError)
			return
		}
		if createdBy.Valid {
			cid := int(createdBy.Int32)
			survey.CreatedBy = &cid
		} else {
			survey.CreatedBy = nil
		}
		surveys = append(surveys, survey)
	}

	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			http.Error(w, "Неавторизованный доступ: нет токена", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка при получении куки", http.StatusBadRequest)
		return
	}

	tokenStr := c.Value
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			http.Error(w, "Неверная подпись токена", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Невалидный токен", http.StatusBadRequest)
		return
	}
	if !token.Valid {
		http.Error(w, "Невалидный токен", http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(10 * time.Second) // Токен действует 24 часа
	claims = &Claims{
		AdminID: claims.AdminID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JwtKey)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Установка куки
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expirationTime,
		HttpOnly: true,
		Secure:   false, // Установите true при использовании HTTPS
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	})

	// Отправка ответа в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveys)
}

// EditSurveyHandler редактирует существующий опрос
func EditSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение SurveyID из URL (/api/surveys/{id})
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
		return
	}
	surveyID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
		return
	}

	// Получение AdminID из контекста
	adminID, ok := GetAdminID(r.Context())
	if !ok {
		http.Error(w, "Не удалось получить AdminID", http.StatusInternalServerError)
		return
	}

	var survey models.Survey
	err = json.NewDecoder(r.Body).Decode(&survey)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if survey.Title == "" {
		http.Error(w, "Название опроса обязательно", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Обновление опроса
	_, err = tx.Exec(
		"UPDATE Surveys SET Title=$1, Description=$2, CreatedBy=$3 WHERE SurveyID=$4",
		survey.Title, survey.Description, adminID, surveyID)
	if err != nil {
		http.Error(w, "Ошибка при обновлении опроса", http.StatusInternalServerError)
		return
	}

	// Удаление существующих вопросов и вариантов ответов
	_, err = tx.Exec(
		"DELETE FROM AnswerOptions WHERE QuestionID IN (SELECT QuestionID FROM Questions WHERE SurveyID=$1)",
		surveyID)
	if err != nil {
		http.Error(w, "Ошибка при удалении вариантов ответов", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(
		"DELETE FROM Questions WHERE SurveyID=$1",
		surveyID)
	if err != nil {
		http.Error(w, "Ошибка при удалении вопросов", http.StatusInternalServerError)
		return
	}

	// Добавление новых вопросов и вариантов ответов
	for _, q := range survey.Questions {
		var questionID int
		err = tx.QueryRow(
			"INSERT INTO Questions (SurveyID, QuestionText, QuestionType) VALUES ($1, $2, $3) RETURNING QuestionID",
			surveyID, q.QuestionText, q.QuestionType).Scan(&questionID)
		if err != nil {
			http.Error(w, "Ошибка при добавлении вопроса", http.StatusInternalServerError)
			return
		}

		if q.QuestionType == "single_choice" || q.QuestionType == "multiple_choice" {
			for _, opt := range q.Options {
				_, err = tx.Exec(
					"INSERT INTO AnswerOptions (QuestionID, OptionText) VALUES ($1, $2)",
					questionID, opt.OptionText)
				if err != nil {
					http.Error(w, "Ошибка при добавлении варианта ответа", http.StatusInternalServerError)
					return
				}
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		http.Error(w, "Ошибка при сохранении изменений", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Опрос успешно обновлён",
	})
}

// GetSurveyHandler возвращает информацию об опросе по его ID
func GetSurveyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		SubmitSurveyHandler(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение SurveyID из URL (/api/surveys/{id})
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 {
		http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
		return
	}
	surveyID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
		return
	}

	var survey models.Survey
	err = db.DB.QueryRow(
		"SELECT SurveyID, Title, Description, CreatedAt, CreatedBy FROM Surveys WHERE SurveyID=$1",
		surveyID).Scan(&survey.SurveyID, &survey.Title, &survey.Description, &survey.CreatedAt, &survey.CreatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Опрос не найден", http.StatusNotFound)
			return
		}
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Получение вопросов
	rows, err := db.DB.Query(
		"SELECT QuestionID, QuestionText, QuestionType FROM Questions WHERE SurveyID=$1",
		surveyID)
	if err != nil {
		http.Error(w, "Ошибка при получении вопросов", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var q models.Question
		err := rows.Scan(&q.QuestionID, &q.QuestionText, &q.QuestionType)
		if err != nil {
			http.Error(w, "Ошибка при обработке вопросов", http.StatusInternalServerError)
			return
		}

		if q.QuestionType == "single_choice" || q.QuestionType == "multiple_choice" {
			optRows, err := db.DB.Query(
				"SELECT OptionID, OptionText FROM AnswerOptions WHERE QuestionID=$1",
				q.QuestionID)
			if err != nil {
				http.Error(w, "Ошибка при получении вариантов ответов", http.StatusInternalServerError)
				return
			}

			for optRows.Next() {
				var opt models.Option
				err := optRows.Scan(&opt.OptionID, &opt.OptionText)
				if err != nil {
					http.Error(w, "Ошибка при обработке вариантов ответов", http.StatusInternalServerError)
					optRows.Close()
					return
				}
				q.Options = append(q.Options, opt)
			}
			optRows.Close()
		}

		survey.Questions = append(survey.Questions, q)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(survey)
}
