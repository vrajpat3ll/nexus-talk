package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Minimal identity service backed by Postgres to allow two real users on one laptop.

var (
	db *sql.DB

	userPresence = map[string]struct {
		Online   bool
		LastSeen time.Time
	}{}
	presenceMux sync.Mutex
	userPrivacy = map[string]map[string]bool{}
	privacyMux  sync.Mutex
)

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}

func withCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}

func ensureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT now()
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	conn := os.Getenv("DB_CONN")
	if conn != "" {
		var err error
		db, err = sql.Open("postgres", conn)
		if err != nil {
			log.Fatalf("identity: failed to open DB: %v", err)
		}
		if err = db.Ping(); err != nil {
			log.Fatalf("identity: DB not reachable: %v", err)
		}
		if err := ensureSchema(db); err != nil {
			log.Fatalf("identity: ensure schema failed: %v", err)
		}
		log.Println("Identity connected to Postgres")
	} else {
		log.Println("WARNING: DB_CONN not set; identity requires Postgres per design doc")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/register", registerHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/profile/", profileHandler)
	// simple lookup by username: /profile/by-username?u=
	mux.HandleFunc("/profile/by-username", profileByUsernameHandler)
	mux.HandleFunc("/presence", presenceHandler)
	mux.HandleFunc("/privacy", privacyHandler)
	log.Println("Identity service listening on :8081 (CORS enabled)")
	log.Fatal(http.ListenAndServe(":8081", withCORS(mux)))
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	if body.Username == "" || body.Password == "" {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	id := uuid.New().String()
	if _, err := db.Exec(`INSERT INTO users (id,username,email,password_hash) VALUES ($1,$2,$3,$4)`, id, body.Username, body.Email, string(hash)); err != nil {
		http.Error(w, "username taken or DB error", http.StatusConflict)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id, "username": body.Username})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	var id string
	var hash string
	if err := db.QueryRow(`SELECT id, password_hash FROM users WHERE username=$1`, body.Username).Scan(&id, &hash); err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(body.Password)) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	// For prototype, return a simple token (not a real JWT) and the user object
	resp := map[string]interface{}{
		"token": "dev-token-" + id,
		"user":  map[string]string{"id": id, "username": body.Username},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func profileHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/profile/")
	var u User
	if err := db.QueryRow(`SELECT id, username, COALESCE(email,'') FROM users WHERE id=$1`, id).Scan(&u.ID, &u.Username, &u.Email); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

func profileByUsernameHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("u")
	if username == "" {
		http.Error(w, "missing u", http.StatusBadRequest)
		return
	}
	var u User
	if err := db.QueryRow(`SELECT id, username, COALESCE(email,'') FROM users WHERE username=$1`, username).Scan(&u.ID, &u.Username, &u.Email); err != nil {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(u)
}

// /presence and /privacy kept as in-memory for demo; can be moved to Postgres later
func presenceHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := r.URL.Query().Get("user_id")
		presenceMux.Lock()
		p, userOk := userPresence[user]
		presenceMux.Unlock()
		j := make(map[string]interface{})
		j["online"] = userOk && p.Online
		j["last_seen"] = p.LastSeen.Format(time.RFC3339)
		json.NewEncoder(w).Encode(j)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			UserID string
			Online bool
		}
		json.NewDecoder(r.Body).Decode(&req)
		presenceMux.Lock()
		userPresence[req.UserID] = struct {
			Online   bool
			LastSeen time.Time
		}{Online: req.Online, LastSeen: time.Now()}
		presenceMux.Unlock()
		w.Write([]byte(`{"ok":true}`))
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}

func privacyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		user := r.URL.Query().Get("user_id")
		privacyMux.Lock()
		settings, ok := userPrivacy[user]
		privacyMux.Unlock()
		if !ok {
			settings = map[string]bool{"profile_visible": true, "last_seen_visible": true, "read_receipts": true}
		}
		json.NewEncoder(w).Encode(settings)
		return
	} else if r.Method == http.MethodPost {
		var req struct {
			UserID   string
			Settings map[string]bool
		}
		json.NewDecoder(r.Body).Decode(&req)
		privacyMux.Lock()
		userPrivacy[req.UserID] = req.Settings
		privacyMux.Unlock()
		w.Write([]byte(`{"ok":true}`))
		return
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
