package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/transactions", transactionsHandler)
	mux.HandleFunc("/subscribe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "subscription": "ok"}`)
	})
	mux.HandleFunc("/unsubscribe", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	log.Println("Payments service listening on :8088")
	log.Fatal(http.ListenAndServe(":8088", mux))
}

func transactionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "tx_id": "tx-stub"}`)
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"transactions": []string{"demo payment 1", "demo payment 2"},
	})
}
