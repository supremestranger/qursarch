package server

import (
	"fmt"
	"log"
	"net/http"
)

func Log() {
	log.Println("server working")
}

func onGet(rw http.ResponseWriter, request *http.Request) {
	rw.Write([]byte("get"))
}

func InitRouting() {
	RegisterOnGet("/get", onGet)
}

func RegisterOnGet(url string, handler http.HandlerFunc) {
	pattern := fmt.Sprintf("GET %s", url)
	http.HandleFunc(pattern, handler)
}

func RegisterOnPost(url string, handler http.HandlerFunc) {
	pattern := fmt.Sprintf("POST %s", url)
	http.HandleFunc(pattern, handler)
}
