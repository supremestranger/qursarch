// handlers/analytics.go
package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "survey-platform-server/db"
    "survey-platform-server/models"
)

func HeatmapHandler(w http.ResponseWriter, r *http.Request, surveyID int) {
    // Получение вопросов и подсчёт ответов
    rows, err := db.DB.Query(`
        SELECT q.QuestionText, COUNT(a.AnswerID)
        FROM Questions q
        LEFT JOIN Answers a ON q.QuestionID = a.QuestionID
        WHERE q.SurveyID = $1
        GROUP BY q.QuestionText
    `, surveyID)
    if err != nil {
        http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    heatmap := make(models.HeatmapData)
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

func ResponseDistributionHandler(w http.ResponseWriter, r *http.Request, surveyID int) {
    // Получение распределения ответов по вариантам
    rows, err := db.DB.Query(`
        SELECT a.AnswerText, COUNT(a.AnswerID)
        FROM Answers a
        JOIN Questions q ON a.QuestionID = q.QuestionID
        WHERE q.SurveyID = $1 AND a.AnswerText IS NOT NULL
        GROUP BY a.AnswerText
    `, surveyID)
    if err != nil {
        http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    data := models.ResponseDistributionData{
        Labels: []string{},
        Counts: []int{},
    }

    for rows.Next() {
        var answer string
        var count int
        if err := rows.Scan(&answer, &count); err != nil {
            http.Error(w, "Ошибка при обработке данных", http.StatusInternalServerError)
            return
        }
        data.Labels = append(data.Labels, answer)
        data.Counts = append(data.Counts, count)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func AverageScoreHandler(w http.ResponseWriter, r *http.Request, surveyID int) {
    // Предполагается, что ответы имеют баллы (например, для оценочных вопросов)
    var average float64
    err := db.DB.QueryRow(`
        SELECT AVG(a.Score)
        FROM Answers a
        JOIN Questions q ON a.QuestionID = q.QuestionID
        WHERE q.SurveyID = $1 AND a.Score IS NOT NULL
    `, surveyID).Scan(&average)
    if err != nil {
        if err == sql.ErrNoRows {
            average = 0
        } else {
            http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
            return
        }
    }

    data := models.AverageScoreData{
        AverageScore: average,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

func CompletionRateHandler(w http.ResponseWriter, r *http.Request, surveyID int) {
    // Подсчёт общего количества пользователей, начавших опрос
    var totalUsers int
    err := db.DB.QueryRow(`
        SELECT COUNT(DISTINCT UserID)
        FROM Answers a
        JOIN Questions q ON a.QuestionID = q.QuestionID
        WHERE q.SurveyID = $1
    `, surveyID).Scan(&totalUsers)
    if err != nil {
        http.Error(w, "Ошибка при подсчёте пользователей", http.StatusInternalServerError)
        return
    }

    // Подсчёт пользователей, завершивших опрос (ответили на все вопросы)
    var completedUsers int
    err = db.DB.QueryRow(`
        SELECT COUNT(*)
        FROM (
            SELECT a.UserID
            FROM Answers a
            JOIN Questions q ON a.QuestionID = q.QuestionID
            WHERE q.SurveyID = $1
            GROUP BY a.UserID
            HAVING COUNT(a.QuestionID) = (SELECT COUNT(*) FROM Questions WHERE SurveyID = $1)
        ) AS completed
    `, surveyID).Scan(&completedUsers)
    if err != nil {
        if err == sql.ErrNoRows {
            completedUsers = 0
        } else {
            http.Error(w, "Ошибка при подсчёте завершивших пользователей", http.StatusInternalServerError)
            return
        }
    }

    completionRate := 0.0
    if totalUsers > 0 {
        completionRate = (float64(completedUsers) / float64(totalUsers)) * 100
    }

    data := models.CompletionRateData{
        CompletionRate: completionRate,
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(data)
}

// AnalysisHandler направляет запросы к соответствующим функциям анализа
func AnalysisHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение SurveyID из URL (/analytics/{id}?type=...)
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 3 {
        http.Error(w, "Не указан ID опроса", http.StatusBadRequest)
        return
    }
    surveyID, err := strconv.Atoi(parts[2])
    if err != nil {
        http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
        return
    }

    query := r.URL.Query()
    analysisType := query.Get("type")

    switch analysisType {
    case "heatmap":
        HeatmapHandler(w, r, surveyID)
    case "response_distribution":
        ResponseDistributionHandler(w, r, surveyID)
    case "average_score":
        AverageScoreHandler(w, r, surveyID)
    case "completion_rate":
        CompletionRateHandler(w, r, surveyID)
    default:
        http.Error(w, "Неизвестный тип анализа", http.StatusBadRequest)
    }
}
