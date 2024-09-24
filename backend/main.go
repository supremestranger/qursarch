package main

import (
	"backend/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Startup!!!а")
	server.InitRouting()
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Println(err)
	}
	// todo server
}
