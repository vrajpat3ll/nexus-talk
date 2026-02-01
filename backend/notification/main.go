package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			w.WriteHeader(201)
			fmt.Fprintln(w, `{ "id": "notif-stub"}`)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"notifications": []string{"Welcome!", "Ping"},
		})
	})
	mux.HandleFunc("/register_device", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "device": "registered" }`)
	})
	log.Println("Notification service listening on :8087")
	log.Fatal(http.ListenAndServe(":8087", mux))
}
