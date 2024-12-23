// handlers/analytics.go
package handlers

import (
	// "database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"survey-platform-server/db"
	// "survey-platform-server/models"
)

// AnalyticsHandler обрабатывает запросы на получение аналитики по опросу
func AnalyticsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Извлечение SurveyID из URL (/api/surveys/{id}/analytics?type=...)
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 5 || parts[4] != "analytics" {
		http.Error(w, "Неверный URL", http.StatusBadRequest)
		return
	}
	surveyID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
		return
	}

	// Получение типа анализа из параметров запроса
	analysisType := r.URL.Query().Get("type")
	if analysisType == "" {
		http.Error(w, "Не указан тип анализа", http.StatusBadRequest)
		return
	}

	switch analysisType {
	case "heatmap":
		getHeatmap(w, surveyID)
	case "pie_chart":
		getPieChart(w, surveyID)
	default:
		http.Error(w, "Неизвестный тип анализа", http.StatusBadRequest)
	}
}

// getHeatmap возвращает количество ответов по каждому вопросу
func getHeatmap(w http.ResponseWriter, surveyID int) {
	query := `
        SELECT q.QuestionText, COUNT(a.AnswerID)
        FROM Questions q
        LEFT JOIN Answers a ON q.QuestionID = a.QuestionID
        WHERE q.SurveyID = $1
        GROUP BY q.QuestionText
    `
	rows, err := db.DB.Query(query, surveyID)
	if err != nil {
		http.Error(w, "Ошибка при получении данных для тепловой карты", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	heatmap := make(map[string]int)
	for rows.Next() {
		var question string
		var count int
		if err := rows.Scan(&question, &count); err != nil {
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
			return
		}
		heatmap[question] = count
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(heatmap)
}

// getPieChart возвращает распределение ответов по вариантам для вопросов с выбором
func getPieChart(w http.ResponseWriter, surveyID int) {
	query := `
        SELECT q.QuestionText, ao.OptionText, COUNT(a.AnswerID)
        FROM Questions q
        JOIN AnswerOptions ao ON q.QuestionID = ao.QuestionID
        LEFT JOIN Answers a ON q.QuestionID = a.QuestionID AND a.SelectedOptions LIKE '%' || ao.OptionText || '%'
        WHERE q.SurveyID = $1 AND q.QuestionType IN ('single_choice', 'multiple_choice')
        GROUP BY q.QuestionText, ao.OptionText
    `
	rows, err := db.DB.Query(query, surveyID)
	if err != nil {
		http.Error(w, "Ошибка при получении данных для пай-чарта", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Структура для хранения данных
	type PieData struct {
		Question string         `json:"question"`
		Options  map[string]int `json:"options"`
	}

	pieDataMap := make(map[string]map[string]int)

	for rows.Next() {
		var question, option string
		var count int
		if err := rows.Scan(&question, &option, &count); err != nil {
			http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
			return
		}

		if _, exists := pieDataMap[question]; !exists {
			pieDataMap[question] = make(map[string]int)
		}
		pieDataMap[question][option] = count
	}

	// Преобразование в слайс
	var pieDataSlice []PieData
	for q, opts := range pieDataMap {
		pieDataSlice = append(pieDataSlice, PieData{
			Question: q,
			Options:  opts,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pieDataSlice)
}
