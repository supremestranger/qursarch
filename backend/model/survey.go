package model

import (
	"backend/surveys"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
)

const SURVEY_ROOT = "/surveys"

type NewSurveyRequest struct {
	Title     string             `json:"title"`
	Questions []surveys.Question `json:"questions"`
}

func RegisterSurveyModels() {
	utils.RegisterOnGet(SURVEY_ROOT, onSurveysGet)
	utils.RegisterOnGet(SURVEY_ROOT+"/{id}", onSurveyGet)
	utils.RegisterOnPost(SURVEY_ROOT, onSurveysPost)
}

func onSurveysGet(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw, "*")
	surveys, err := surveys.GetSurveys()
	if err != nil {
		log.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	res, err := json.Marshal(surveys)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	rw.Write(res)
}

func onSurveyGet(rw http.ResponseWriter, req *http.Request) {
	utils.EnableCors(rw, "*")
	id := req.PathValue("id")
	survey, err := surveys.GetSurveyById(id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	surveyJson, err := json.Marshal(survey)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	rw.Write(surveyJson)
}

func onSurveysPost(rw http.ResponseWriter, req *http.Request) {
	log.Println("hello 123123")
	utils.EnableCors(rw, "http://localhost:3000")
	rw.Header().Add("Access-Control-Allow-Credentials", "true")
	ok, user := CheckAuth(rw, req)
	if !ok {
		http.Error(rw, "вы не авторизованы", http.StatusBadRequest)
		return
	}

	var newSurveyReq NewSurveyRequest
	err := json.NewDecoder(req.Body).Decode(&newSurveyReq)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	if len(newSurveyReq.Questions) == 0 || newSurveyReq.Questions == nil {
		http.Error(rw, "no questions", http.StatusBadRequest)
		return
	}

	for _, v := range newSurveyReq.Questions {
		if !TypeIsCorrect(v) {
			http.Error(rw, "no answers", http.StatusBadRequest)
			return
		}
	}

	json, err := json.Marshal(&newSurveyReq.Questions)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	err = surveys.CreateSurvey(string(json), newSurveyReq.Title, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
}

func TypeIsCorrect(q surveys.Question) bool {
	if q.Type != surveys.S_TEXT_ANSWER_TYPE && q.Type != surveys.L_TEXT_ANSWER_TYPE && q.Type != surveys.SINGLE_ANSWER_TYPE && q.Type != surveys.MULTIPLE_ANSWERS_TYPE {
		return false
	}

	noAnswers := q.Answers == nil || len(q.Answers) == 0

	if (q.Type == surveys.SINGLE_ANSWER_TYPE || q.Title == surveys.MULTIPLE_ANSWERS_TYPE) && noAnswers {
		return false
	}

	return true
}
