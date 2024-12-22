package models

type User struct {
    UserID   int    `json:"user_id"`
    Login    string `json:"login"`
    Password string `json:"-"`
}

type Survey struct {
    SurveyID    int       `json:"survey_id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    CreatedAt   string    `json:"created_at"`
    CreatedBy   int       `json:"created_by"`
    Questions   []Question `json:"questions"`
}

type Question struct {
    QuestionID   int      `json:"question_id"`
    QuestionText string   `json:"question_text"`
    QuestionType string   `json:"question_type"`
    Options      []string `json:"options,omitempty"`
}
type HeatmapData map[string]int

type ResponseDistributionData struct {
    Labels []string `json:"labels"`
    Counts []int    `json:"counts"`
}

type AverageScoreData struct {
    AverageScore float64 `json:"average_score"`
}

type CompletionRateData struct {
    CompletionRate float64 `json:"completion_rate"`
}