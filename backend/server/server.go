package server

import "log"

type Server struct {
	Addr string
}

func Something() {
	log.Println("server working")
}
