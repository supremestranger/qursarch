package main

import (
	"backend/server"
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Startup!!!")
	server.InitRouting()
	server.InitDb()
	err := http.ListenAndServe(":3002", nil)
	if err != nil {
		log.Println(err)
	}
	// todo server
}
