package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "id": "report-stub" }`)
	})
	mux.HandleFunc("/block", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "blocked": true }`)
	})
	mux.HandleFunc("/unblock", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(204)
	})
	mux.HandleFunc("/strikes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(201)
		fmt.Fprintln(w, `{ "strike": "ok" }`)
	})
	log.Println("Moderation service listening on :8086")
	log.Fatal(http.ListenAndServe(":8086", mux))
}
