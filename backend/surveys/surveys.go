package surveys

const (
	SINGLE_ANSWER_TYPE    = "single_answer"
	MULTIPLE_ANSWERS_TYPE = "multiple_answers"
	S_TEXT_ANSWER_TYPE    = "s_text_answer"
	L_TEXT_ANSWER_TYPE    = "l_text_answer"
)

type Question struct {
	Title   string
	Type    string
	Answers []string
}

type Survey struct {
	Questions []Question
}

func GetSurveyById(rawId string) *Survey {
	return &Survey{}
}
