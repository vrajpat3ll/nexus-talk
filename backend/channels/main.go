package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
)

type Channel struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Owner   string   `json:"owner"`
	Members []string `json:"members"`
}

var (
	channels = make(map[string]*Channel)
	chMux    sync.Mutex
	chID     int
)

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create", createChannel)
	mux.HandleFunc("/join", joinChannel)
	mux.HandleFunc("/list", listChannels)
	log.Println("Channels service w/ in-memory demo, CORS enabled (port 8083)")
	log.Fatal(http.ListenAndServe(":8083", withCORS(mux)))
}

func createChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req struct{ Name, Owner string }
	json.NewDecoder(r.Body).Decode(&req)
	chMux.Lock()
	chID++
	id := "c" + string(rune(chID+'A'))
	channels[id] = &Channel{ID: id, Name: req.Name, Owner: req.Owner, Members: []string{req.Owner}}
	chMux.Unlock()
	json.NewEncoder(w).Encode(channels[id])
}
func joinChannel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var req struct{ ChannelID, User string }
	json.NewDecoder(r.Body).Decode(&req)
	chMux.Lock()
	if ch, ok := channels[req.ChannelID]; ok {
		for _, u := range ch.Members {
			if u == req.User {
				return
			}
		}
		ch.Members = append(ch.Members, req.User)
	}
	chMux.Unlock()
	w.Write([]byte(`{"ok":true}`))
}
func listChannels(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	user := r.URL.Query().Get("user_id")
	var result []*Channel
	chMux.Lock()
	for _, ch := range channels {
		for _, u := range ch.Members {
			if u == user {
				result = append(result, ch)
			}
		}
	}
	chMux.Unlock()
	json.NewEncoder(w).Encode(result)
}
