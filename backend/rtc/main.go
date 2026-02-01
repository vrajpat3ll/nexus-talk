package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/calls", callsHandler)
	mux.HandleFunc("/calls/", callSubHandler)
	log.Println("RTC service listening on :8085")
	log.Fatal(http.ListenAndServe(":8085", mux))
}

func callsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "call_id": "call-xxx" }`)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"calls": []string{"demo call 1", "demo call 2"},
	})
}

func callSubHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/calls/"), "/")
	if len(parts) > 1 && parts[1] == "join" {
		w.WriteHeader(204)
		return
	}
	if len(parts) > 1 && parts[1] == "participants" {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"participants": []string{"user1", "user2"},
		})
		return
	}
	if len(parts) > 1 && parts[1] == "end" {
		w.WriteHeader(204)
		return
	}
}
