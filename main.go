package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/health", health)
	http.HandleFunc("/search", search)

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