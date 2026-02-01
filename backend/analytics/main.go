package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "id": "event-stub"}`)
	})
	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"metrics": []string{"m1", "m2"},
		})
	})
	mux.HandleFunc("/counters", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"counters": []string{"c1: 5", "c2: 10"},
		})
	})
	log.Println("Analytics service listening on :8089")
	log.Fatal(http.ListenAndServe(":8089", mux))
}
