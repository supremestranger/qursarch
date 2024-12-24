package models

import "time"

type Admin struct {
    AdminID  int    `json:"admin_id"`
    Login    string `json:"login"`
    Password string `json:"-"`
}

type Survey struct {
    SurveyID    int       `json:"survey_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   time.Time `json:"created_at"`
    CreatedBy   *int       `json:"created_by"`
    Questions   []Question `json:"questions"`
}

type Question struct {
    QuestionID   int      `json:"question_id"`
    SurveyID     int      `json:"survey_id"`
    QuestionText string   `json:"question_text"`
    QuestionType string   `json:"question_type"` // 'single_choice', 'multiple_choice', 'free_text'
    Options      []Option `json:"options,omitempty"`
}

type Option struct {
    OptionID   int    `json:"option_id"`
    QuestionID int    `json:"question_id"`
    OptionText string `json:"option_text"`
}

type SurveyResult struct {
    ResultID    int       `json:"result_id"`
    SurveyID    int       `json:"survey_id"`
    UserID      string    `json:"user_id"` // Рандомно сгенерированный ID пользователя
    CompletedAt time.Time `json:"completed_at"`
    Answers     []Answer  `json:"answers"`
}

type Answer struct {
    AnswerID      int    `json:"answer_id"`
    ResultID      int    `json:"result_id"`
    QuestionID    int    `json:"question_id"`
    AnswerText    string `json:"answer_text,omitempty"`      // Для типа 'free_text'
    SelectedOptions string `json:"selected_options,omitempty"` // Для типов 'single_choice' и 'multiple_choice'
}