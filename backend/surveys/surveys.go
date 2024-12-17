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
	Id        int
	Title     string
	CreatorId int
	Questions string
}

func GetSurveyById(rawId string) (*Survey, error) {
	var survey Survey
	rows, err := db.DB.Query("SELECT * FROM SURVEYS WHERE SURVEYS.ID = $1", rawId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		if err := rows.Scan(&survey.Id, &survey.Title, &survey.Questions, &survey.CreatorId); err != nil {
			return nil, err
		}
		break
	}
	return &survey, nil
}

func GetSurveys() ([]Survey, error) {
	rows, err := db.DB.Query("SELECT ID, Title, Questions, Creator FROM Surveys")
	if err != nil {
		return nil, err
	}
	var surveys []Survey
	for rows.Next() {
		var survey Survey
		if err := rows.Scan(&survey.Id, &survey.Title, &survey.Questions, &survey.CreatorId); err != nil {
			return nil, err
		}
		surveys = append(surveys, survey)
	}

	return surveys, err
}

func CreateSurvey(questionsJson string, title string, user string) error {
	row := db.DB.QueryRow("SELECT Accounts.ID FROM ACCOUNTS Where Accounts.Username = $1", user)

	var id int
	row.Scan(&id)

	res, err := db.DB.Exec("INSERT INTO Surveys (Questions, Title, Creator) values ($1, $2, $3)", questionsJson, title, id)
	if err != nil {
		return err
	}

	log.Println(questionsJson)

	log.Println(res)
	return nil
}
