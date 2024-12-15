package utils

import (
	"fmt"
	"net/http"
)

const (
	API_VERSION = "/v1"
)

func EnableCors(w http.ResponseWriter, val string) {
	w.Header().Add("Access-Control-Allow-Origin", val)
	w.Header().Add("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
}

func RegisterOnGet(url string, handler http.HandlerFunc) {
	pattern := fmt.Sprintf("GET %s", API_VERSION+url)
	http.HandleFunc(pattern, handler)
}

func RegisterOnPost(url string, handler http.HandlerFunc) {
	pattern := fmt.Sprintf("POST %s", API_VERSION+url)
	http.HandleFunc(pattern, handler)
}
