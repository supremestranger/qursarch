package server

import (
	"backend/model"
	"log"
)

func Log() {
	log.Println("server working")
}

func InitRouting() {
	model.RegisterAccountModels()
	model.RegisterSurveyModels()
}
