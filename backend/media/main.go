package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
)

type Media struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Owner string `json:"owner"`
}

var (
	medias   []Media
	mediaMux sync.Mutex
	mediaID  int
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
	mux.HandleFunc("/upload", uploadHandler)
	mux.HandleFunc("/download", downloadHandler)
	log.Println("Media service running (CORS enabled), port 8084")
	log.Fatal(http.ListenAndServe(":8084", withCORS(mux)))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	// Accepts multipart: file, owner
	r.ParseMultipartForm(12 << 20)
	owner := r.FormValue("owner")
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(400)
		fmt.Fprint(w, "file required")
		return
	}
	_ = file
	mediaMux.Lock()
	mediaID++
	medias = append(medias, Media{ID: mediaID, Name: fileHeader.Filename, Owner: owner})
	mediaMux.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"id": mediaID, "name": fileHeader.Filename})
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)
	mediaMux.Lock()
	for _, m := range medias {
		if m.ID == id {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", "attachment; filename="+m.Name)
			io.WriteString(w, "This is a dummy file for id: "+idStr)
			mediaMux.Unlock()
			return
		}
	}
	mediaMux.Unlock()
	w.WriteHeader(404)
}
