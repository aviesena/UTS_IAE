package main

import (
	"encoding/json"
	"log"
	"net/http"
)

const servicePort = "3000"

func main() {
	http.HandleFunc("/", getData)
	log.Fatal(http.ListenAndServe(":3002", nil))
}

func getData(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get("http://localhost:" + servicePort)
	if err != nil {
		log.Println("Error fetching data from service:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println("Error decoding response:", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
