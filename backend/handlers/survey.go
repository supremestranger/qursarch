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
)

func CreateSurveyHandler(w http.ResponseWriter, r *http.Request) {
    // Получение UserID из контекста
    userID, ok := GetUserID(r.Context())
    if !ok {
        http.Error(w, "Не удалось получить UserID", http.StatusInternalServerError)
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
        "INSERT INTO surveys (title, description, createdby) VALUES ($1, $2, $3) RETURNING surveyid",
        survey.Title, survey.Description, userID).Scan(&surveyID)
    if err != nil {
        http.Error(w, "Ошибка при создании опроса", http.StatusInternalServerError)
        return
    }

    for _, q := range survey.Questions {
        var questionID int
        err = tx.QueryRow(
            "INSERT INTO questions (surveyid, questiontext, questiontype) VALUES ($1, $2, $3) RETURNING questionid",
            surveyID, q.QuestionText, q.QuestionType).Scan(&questionID)
        if err != nil {
            http.Error(w, "Ошибка при добавлении вопроса", http.StatusInternalServerError)
            return
        }

        if q.QuestionType == "single_choice" || q.QuestionType == "multiple_choice" {
            for _, opt := range q.Options {
                _, err = tx.Exec(
                    "INSERT INTO answeroptions (questionid, optiontext) VALUES ($1, $2)",
                    questionID, opt)
                if err != nil {
                    http.Error(w, "Ошибка при добавлении вариантов ответа", http.StatusInternalServerError)
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

func GetSurveyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение SurveyID из URL (/survey/{id})
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

    var survey models.Survey
    err = db.DB.QueryRow(
        "SELECT surveyid, title, description, createdat, createdby FROM surveys WHERE surveyid=$1",
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
        "SELECT questionid, questiontext, questiontype FROM questions WHERE surveyid=$1",
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
                "SELECT optiontext FROM answeroptions WHERE questionid=$1",
                q.QuestionID)
            if err != nil {
                http.Error(w, "Ошибка при получении вариантов ответов", http.StatusInternalServerError)
                return
            }

            for optRows.Next() {
                var opt string
                err := optRows.Scan(&opt)
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

func EditSurveyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPut {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение SurveyID из URL (/survey/edit/{id})
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

    // Получение UserID из контекста (можно использовать для проверки прав)
    userID, ok := GetUserID(r.Context())
    if !ok {
        http.Error(w, "Не удалось получить UserID", http.StatusInternalServerError)
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

    // Обновление заголовка и описания опроса
    _, err = tx.Exec(
        "UPDATE surveys SET title=$1, description=$2 WHERE surveyid=$3",
        survey.Title, survey.Description, surveyID)
    if err != nil {
        http.Error(w, "Ошибка при обновлении опроса", http.StatusInternalServerError)
        return
    }

    // Удаление существующих вопросов и вариантов ответов
    _, err = tx.Exec(
        "DELETE FROM answeroptions WHERE questionid IN (SELECT questionid FROM questions WHERE surveyid=$1)",
        surveyID)
    if err != nil {
        http.Error(w, "Ошибка при удалении вариантов ответов", http.StatusInternalServerError)
        return
    }

    _, err = tx.Exec(
        "DELETE FROM questions WHERE surveyid=$1",
        surveyID)
    if err != nil {
        http.Error(w, "Ошибка при удалении вопросов", http.StatusInternalServerError)
        return
    }

    // Добавление новых вопросов и вариантов ответов
    for _, q := range survey.Questions {
        var questionID int
        err = tx.QueryRow(
            "INSERT INTO questions (surveyid, questiontext, questiontype) VALUES ($1, $2, $3) RETURNING questionid",
            surveyID, q.QuestionText, q.QuestionType).Scan(&questionID)
        if err != nil {
            http.Error(w, "Ошибка при добавлении вопроса", http.StatusInternalServerError)
            return
        }

        if q.QuestionType == "single_choice" || q.QuestionType == "multiple_choice" {
            for _, opt := range q.Options {
                _, err = tx.Exec(
                    "INSERT INTO answeroptions (questionid, optiontext) VALUES ($1, $2)",
                    questionID, opt)
                if err != nil {
                    http.Error(w, "Ошибка при добавлении вариантов ответа", http.StatusInternalServerError)
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

func SubmitSurveyHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
        return
    }

    // Извлечение SurveyID из URL (/survey/{id}/submit)
    parts := strings.Split(r.URL.Path, "/")
    if len(parts) < 4 || parts[3] != "submit" {
        http.Error(w, "Неверный URL", http.StatusBadRequest)
        return
    }
    surveyID, err := strconv.Atoi(parts[2])
    if err != nil {
        http.Error(w, "Неверный ID опроса", http.StatusBadRequest)
        return
    }

    var submission struct {
        UserID  int                    `json:"user_id"`
        Answers map[string]interface{} `json:"answers"`
    }

    err = json.NewDecoder(r.Body).Decode(&submission)
    if err != nil {
        http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
        return
    }

    if submission.UserID == 0 || len(submission.Answers) == 0 {
        http.Error(w, "Необходимо указать user_id и ответы", http.StatusBadRequest)
        return
    }

    // Получение UserID из контекста и проверка соответствия
    authUserID, ok := GetUserID(r.Context())
    if !ok {
        http.Error(w, "Не удалось получить UserID из токена", http.StatusInternalServerError)
        return
    }

    if authUserID != submission.UserID {
        http.Error(w, "Несоответствие UserID", http.StatusUnauthorized)
        return
    }

    tx, err := db.DB.Begin()
    if err != nil {
        http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
        return
    }
    defer tx.Rollback()

    for qIDStr, answer := range submission.Answers {
        qID, err := strconv.Atoi(qIDStr)
        if err != nil {
            http.Error(w, "Неверный формат question_id", http.StatusBadRequest)
            return
        }

        switch ans := answer.(type) {
        case map[string]interface{}:
            if selected, ok := ans["selected_options"]; ok {
                // Обработка выбранных вариантов
                selectedOptions, ok := selected.([]interface{})
                if !ok {
                    http.Error(w, "Неверный формат выбранных опций", http.StatusBadRequest)
                    return
                }

                for _, opt := range selectedOptions {
                    optStr, ok := opt.(string)
                    if !ok {
                        http.Error(w, "Неверный формат опции", http.StatusBadRequest)
                        return
                    }

                    _, err = tx.Exec(
                        "INSERT INTO answers (userid, questionid, answertext) VALUES ($1, $2, $3)",
                        submission.UserID, qID, optStr)
                    if err != nil {
                        http.Error(w, "Ошибка при сохранении ответа", http.StatusInternalServerError)
                        return
                    }
                }
            }
            if text, ok := ans["answer_text"]; ok {
                answerText, ok := text.(string)
                if !ok {
                    http.Error(w, "Неверный формат текста ответа", http.StatusBadRequest)
                    return
                }

                _, err = tx.Exec(
                    "INSERT INTO answers (userid, questionid, answertext) VALUES ($1, $2, $3)",
                    submission.UserID, qID, answerText)
                if err != nil {
                    http.Error(w, "Ошибка при сохранении ответа", http.StatusInternalServerError)
                    return
                }
            }
        default:
            http.Error(w, "Неверный формат ответа", http.StatusBadRequest)
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
