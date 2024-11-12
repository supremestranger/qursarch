package surveys

import (
	"backend/db"
	"log"
)

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

func CreateSurvey(questionsJson string, user string) error {
	row := db.DB.QueryRow("SELECT Accounts.ID FROM ACCOUNTS Where Accounts.Username = $1", user)

	var id int
	row.Scan(&id)

	res, err := db.DB.Exec("INSERT INTO Surveys (Questions, Creator) values ($1, $2)", questionsJson, id)
	if err != nil {
		return err
	}

	log.Println(questionsJson)

	log.Println(res)
	return nil
}
