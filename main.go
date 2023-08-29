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
	fmt.Println("health is called")
	response := HealthResponse{
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
	Number  string `json:"number"`
}

func search(w http.ResponseWriter, r *http.Request) {

	number := r.URL.Query().Get("number")

	response := SearchResponse{
		Message: "search is called !",
		Number:  number,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(response)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type DevicesResponse struct {
	Names []string `json:"names"`
}

func devices(w http.ResponseWriter, r *http.Request) {

	deviceName := r.URL.Query().Get("name")

	fmt.Println("Endpoint /devices is called", deviceName)

	// open database connection
	// 		account: vince
	// 		password: 1234
	// 		host: 127.0.0.1
	// 		port: 3306
	// 		schema: vince_test
	db, err := sql.Open("mysql", "vince:1234@tcp(127.0.0.1:3306)/vince_test")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	// execute sql
	rows, err := db.Query("SELECT name FROM DEVICE WHERE name like '%" + deviceName + "%';")

	if err != nil {
		fmt.Println("Error executing SQL query:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// generate http response
	var names []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			fmt.Println("Error scanning row:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		names = append(names, name)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error with rows:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(names) == 0 {
		fmt.Println("No data found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	response := DevicesResponse{
		Names: names,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
	}
}
