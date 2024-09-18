package main

import (
	"backend/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("hello world")
	server.InitRouting()
	err := http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Println(err)
	}
	// todo server
}
