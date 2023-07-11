package main

import (
	"encoding/json"
	"fmt"
	"net/http"
    
	// DB
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/search", search)
	http.HandleFunc("/devices", devices)

	err := http.ListenAndServe(":80", nil)

	if err != nil {
		fmt.Println("server start fail")
	}
}

type HealthResponse struct {
	Message string `json:"message"`
}

func health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse {
		Message: "Service is healthy !",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type SearchResponse struct {
	Message string `json:"message"`
	Number string `json:"number"`
}

func search(w http.ResponseWriter, r *http.Request) {

	number := r.URL.Query().Get("number")

	response := SearchResponse {
		Message: "search is called !",
		Number: number,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response) 

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func devices(w http.ResponseWriter, r *http.Request) {
	deviceId := r.URL.Query().Get("id")

	// open database connection
	// 		account: vince
	// 		password: 1234
	// 		host: 127.0.0.1
	// 		port: 3306
	// 		schema: vince_test


	// execute sql


    // generate http response

}