// handlers/submit_survey.go
package handlers

import (
    "database/sql"
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "survey-platform-server/db"
    // "survey-platform-server/models"
)

// SubmitSurveyHandler обрабатывает отправку опроса пользователем
func SubmitSurveyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение SurveyID из URL (/api/surveys/{id}/submit)
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 5 || parts[4] != "submit" {
        http.Error(w, "Неверный URL", http.StatusBadRequest)
        return
    }
    surveyID, err := strconv.Atoi(parts[3])
    if err != nil {
        http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
        return
    }

    var submission struct {
        UserID  string                 `json:"user_id"` // Рандомно сгенерированный ID пользователя
        Answers map[string]AnswerInput `json:"answers"`
    }

    err = json.NewDecoder(r.Body).Decode(&submission)
    if err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    if submission.UserID == "" || len(submission.Answers) == 0 {
        http.Error(w, "Необходимо указать user_id и ответы", http.StatusBadRequest)
        return
    }

    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    // Вставка результата опроса
    var resultID int
    err = tx.QueryRow(
        "INSERT INTO SurveyResults (SurveyID, UserID) VALUES ($1, $2) RETURNING ResultID",
        surveyID, submission.UserID).Scan(&resultID)
    if err != nil {
        http.Error(w, "Ошибка при сохранении результата опроса", http.StatusInternalServerError)
        return
    }

    for qIDStr, answer := range submission.Answers {
        qID, err := strconv.Atoi(qIDStr)
        if err != nil {
            http.Error(w, "Неверный формат question_id", http.StatusBadRequest)
            return
        }

        var answerText sql.NullString
        var selectedOptions sql.NullString

        if answer.AnswerText != "" {
            answerText = sql.NullString{String: answer.AnswerText, Valid: true}
        }

        if len(answer.SelectedOptions) > 0 {
            // Преобразование выбранных опций в строку, разделённую запятой
            options := make([]string, len(answer.SelectedOptions))
            for i, opt := range answer.SelectedOptions {
                options[i] = opt.OptionText
            }
            selectedOptions = sql.NullString{String: strings.Join(options, ","), Valid: true}
        }

        _, err = tx.Exec(
            "INSERT INTO Answers (ResultID, QuestionID, AnswerText, SelectedOptions) VALUES ($1, $2, $3, $4)",
            resultID, qID, answerText, selectedOptions)
        if err != nil {
            http.Error(w, "Ошибка при сохранении ответа", http.StatusInternalServerError)
            return
        }
    }

    err = tx.Commit()
    if err != nil {
        http.Error(w, "Ошибка при сохранении ответов", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Ответы успешно сохранены",
    })
}

// AnswerInput представляет структуру ответа, полученного от пользователя
type AnswerInput struct {
    AnswerText      string         `json:"answer_text,omitempty"`
    SelectedOptions []OptionInput  `json:"selected_options,omitempty"`
}

// OptionInput представляет структуру выбранного варианта ответа
type OptionInput struct {
    OptionText string `json:"option_text"`
}
