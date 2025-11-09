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
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq"
)

// Message represents a single chat message persisted in the DB.
type Message struct {
	ID          string                 `json:"id"`
	SenderID    string                 `json:"sender_id"`
	ThreadID    string                 `json:"thread_id"`
	Content     string                 `json:"content"`
	SentAt      time.Time              `json:"sent_at"`
	ContentType string                 `json:"content_type,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

var (
	db         *sql.DB
	memLock    sync.Mutex
	memMsgs    []Message
	memThreads = map[string]string{} // key: direct pair key "u1|u2" -> threadID
	useMemory  bool
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

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // In production, implement proper origin checking
		},
	}
	connManager = NewConnectionManager()
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}

	connManager.AddClient(userID, conn)
	defer func() {
		connManager.RemoveClient(userID, conn)
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket read error: %v", err)
			}
			break
		}
		// Handle incoming WebSocket messages if needed
		log.Printf("Received WebSocket message from user %s: %s", userID, string(msg))
	}
}

func main() {
	conn := os.Getenv("DB_CONN")
	if conn != "" {
		var err error
		db, err = sql.Open("postgres", conn)
		if err != nil {
			log.Printf("DB open error: %v — falling back to memory mode", err)
			useMemory = true
		} else if err = db.Ping(); err != nil {
			log.Printf("DB ping error: %v — falling back to memory mode", err)
			useMemory = true
		} else {
			if err := ensureSchema(db); err != nil {
				log.Printf("ensure schema error: %v — falling back to memory mode", err)
				useMemory = true
			}
		}
	} else {
		useMemory = true
	}

	mux := http.NewServeMux()
	// Create message (1:1 or thread-based)
	mux.HandleFunc("/v1/messages", messagesHandler)
	// Get messages by thread id: /v1/messages/{threadId}
	mux.HandleFunc("/v1/messages/", threadMessagesHandler)
	// List threads for a user: /v1/threads?user_id=xxx
	mux.HandleFunc("/v1/threads", listThreadsHandler)
	// Create or fetch a direct thread for two users
	mux.HandleFunc("/v1/threads/direct", directThreadHandler)
	// WebSocket endpoint for real-time updates
	mux.HandleFunc("/v1/ws", wsHandler)

	// Keep old quick endpoints for local demo compatibility
	mux.HandleFunc("/send", sendDemoHandler)
	mux.HandleFunc("/inbox", inboxDemoHandler)

	log.Printf("Messaging service starting on :8082 (memory=%v)", useMemory)
	log.Fatal(http.ListenAndServe(":8082", withCORS(mux)))
}

// ensureSchema creates minimal tables if they don't exist.
func ensureSchema(db *sql.DB) error {
	stmts := []string{
		`CREATE TABLE IF NOT EXISTS threads (
			id UUID PRIMARY KEY,
			type VARCHAR(16) CHECK (type IN ('direct','group','channel')) DEFAULT 'direct',
			created_at TIMESTAMP DEFAULT now()
		)`,
		`CREATE TABLE IF NOT EXISTS thread_members (
			thread_id UUID REFERENCES threads(id) ON DELETE CASCADE,
			user_id UUID NOT NULL,
			PRIMARY KEY (thread_id, user_id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id UUID PRIMARY KEY,
			sender_id UUID,
			thread_id UUID REFERENCES threads(id),
			content TEXT,
			sent_at TIMESTAMP DEFAULT now(),
			content_type VARCHAR(32),
			is_deleted BOOLEAN DEFAULT FALSE,
			reply_to UUID,
			metadata JSONB
		)`,
		`CREATE TABLE IF NOT EXISTS message_states (
			message_id UUID REFERENCES messages(id),
			user_id UUID,
			state VARCHAR(16),
			timestamp TIMESTAMP DEFAULT now(),
			PRIMARY KEY (message_id, user_id)
		)`,
	}
	for _, s := range stmts {
		if _, err := db.Exec(s); err != nil {
			return err
		}
	}
	return nil
}

// messagesHandler handles POST /v1/messages to create a message. It also
// accepts GET with query params for simple listing (deprecated path).
func messagesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var body struct {
			SenderID    string                 `json:"sender_id"`
			ToID        string                 `json:"to_id,omitempty"`
			ThreadID    string                 `json:"thread_id,omitempty"`
			Content     string                 `json:"content"`
			ContentType string                 `json:"content_type,omitempty"`
			Metadata    map[string]interface{} `json:"metadata,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, "invalid body", http.StatusBadRequest)
			return
		}

		// If no thread supplied, create a direct thread if to_id is present.
		threadID := body.ThreadID
		if threadID == "" && body.ToID != "" {
			// Find or create a direct thread between SenderID and ToID
			u1, u2 := canonicalPair(body.SenderID, body.ToID)
			if useMemory {
				key := u1 + "|" + u2
				memLock.Lock()
				tid, ok := memThreads[key]
				if !ok {
					tid = uuid.New().String()
					memThreads[key] = tid
				}
				memLock.Unlock()
				threadID = tid
			} else {
				tid, err := findOrCreateDirectThread(db, u1, u2)
				if err != nil {
					log.Printf("failed to get/create direct thread: %v", err)
					http.Error(w, "internal", http.StatusInternalServerError)
					return
				}
				threadID = tid
			}
		}

		msg := Message{
			ID:          uuid.New().String(),
			SenderID:    body.SenderID,
			ThreadID:    threadID,
			Content:     body.Content,
			SentAt:      time.Now().UTC(),
			ContentType: body.ContentType,
			Metadata:    body.Metadata,
		}

		if useMemory {
			memLock.Lock()
			memMsgs = append(memMsgs, msg)
			memLock.Unlock()
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(msg)
			log.Printf("[mem] message stored: %s -> thread %s", msg.SenderID, msg.ThreadID)
			return
		}

		// Persist to DB
		metaJSON := "null"
		if msg.Metadata != nil {
			b, _ := json.Marshal(msg.Metadata)
			metaJSON = string(b)
		}
		q := `INSERT INTO messages (id, sender_id, thread_id, content, sent_at, content_type, metadata) VALUES ($1,$2,$3,$4,$5,$6,$7)`
		if _, err := db.Exec(q, msg.ID, msg.SenderID, msg.ThreadID, msg.Content, msg.SentAt, msg.ContentType, metaJSON); err != nil {
			log.Printf("db insert message failed: %v", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}

		// Send response to the sender
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(msg)

		// Notify connected clients in the thread through WebSocket
		notification, _ := json.Marshal(map[string]interface{}{
			"type": "new_message",
			"data": msg,
		})

		// Get thread members
		if !useMemory {
			var members []string
			rows, err := db.Query(`
				SELECT user_id FROM thread_members 
				WHERE thread_id = $1 AND user_id != $2
			`, msg.ThreadID, msg.SenderID)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var userID string
					if err := rows.Scan(&userID); err == nil {
						members = append(members, userID)
					}
				}
				// Update thread members in connection manager
				connManager.UpdateThreadMembers(msg.ThreadID, members)
			}
		}

		// Send real-time notification through WebSocket
		connManager.SendToThread(msg.ThreadID, notification, msg.SenderID)

		// Log the event
		log.Printf("message created: id=%s thread=%s sender=%s", msg.ID, msg.ThreadID, msg.SenderID)

	case http.MethodGet:
		// Support listing by query param thread_id
		thread := r.URL.Query().Get("thread_id")
		if thread == "" {
			http.Error(w, "missing thread_id", http.StatusBadRequest)
			return
		}
		sendMessagesByThread(w, r, thread)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// threadMessagesHandler handles GET /v1/messages/{threadId}
func threadMessagesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	base := "/v1/messages/"
	threadID := strings.TrimPrefix(r.URL.Path, base)
	if threadID == "" {
		http.Error(w, "missing thread id in path", http.StatusBadRequest)
		return
	}
	sendMessagesByThread(w, r, threadID)
}

func sendMessagesByThread(w http.ResponseWriter, r *http.Request, threadID string) {
	// Optional since param as RFC3339 timestamp
	since := r.URL.Query().Get("since")

	if useMemory {
		memLock.Lock()
		out := []Message{}
		for _, m := range memMsgs {
			if m.ThreadID == threadID {
				if since != "" {
					if t, err := time.Parse(time.RFC3339, since); err == nil {
						if m.SentAt.After(t) {
							out = append(out, m)
						}
					}
				} else {
					out = append(out, m)
				}
			}
		}
		memLock.Unlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(out)
		return
	}

	q := `SELECT id, sender_id, thread_id, content, sent_at, content_type, metadata FROM messages WHERE thread_id = $1`
	args := []interface{}{threadID}
	if since != "" {
		if t, err := time.Parse(time.RFC3339, since); err == nil {
			q += " AND sent_at > $2"
			args = append(args, t)
		}
	}
	q += " ORDER BY sent_at ASC"

	rows, err := db.Query(q, args...)
	if err != nil {
		log.Printf("db query failed: %v", err)
		http.Error(w, "internal", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	out := []Message{}
	for rows.Next() {
		var m Message
		var meta sql.NullString
		if err := rows.Scan(&m.ID, &m.SenderID, &m.ThreadID, &m.Content, &m.SentAt, &m.ContentType, &meta); err != nil {
			log.Printf("scan error: %v", err)
			continue
		}
		if meta.Valid {
			var mm map[string]interface{}
			if err := json.Unmarshal([]byte(meta.String), &mm); err == nil {
				m.Metadata = mm
			}
		}
		out = append(out, m)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(out)
}

// directThreadHandler accepts POST {"user_a":"...","user_b":"..."} and returns {thread_id}
func directThreadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	var body struct {
		UserA string `json:"user_a"`
		UserB string `json:"user_b"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserA == "" || body.UserB == "" {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	u1, u2 := canonicalPair(body.UserA, body.UserB)
	var tid string
	if useMemory {
		key := u1 + "|" + u2
		memLock.Lock()
		existing, ok := memThreads[key]
		if !ok {
			existing = uuid.New().String()
			memThreads[key] = existing
		}
		memLock.Unlock()
		tid = existing
	} else {
		var err error
		tid, err = findOrCreateDirectThread(db, u1, u2)
		if err != nil {
			log.Printf("direct thread error: %v", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"thread_id": tid})
}

// canonicalPair sorts two IDs to build a stable key
func canonicalPair(a, b string) (string, string) {
	if a <= b {
		return a, b
	}
	return b, a
}

// findOrCreateDirectThread ensures a single thread for the pair (u1,u2)
func findOrCreateDirectThread(db *sql.DB, u1, u2 string) (string, error) {
	// Try to find existing direct thread with exactly these two members
	var tid string
	q := `SELECT t.id FROM threads t
		  JOIN thread_members m1 ON m1.thread_id=t.id AND m1.user_id=$1
		  JOIN thread_members m2 ON m2.thread_id=t.id AND m2.user_id=$2
		  WHERE t.type='direct'
		  LIMIT 1`
	if err := db.QueryRow(q, u1, u2).Scan(&tid); err == nil {
		return tid, nil
	}
	// Create new
	tid = uuid.New().String()
	if _, err := db.Exec(`INSERT INTO threads (id,type) VALUES ($1,'direct')`, tid); err != nil {
		return "", err
	}
	if _, err := db.Exec(`INSERT INTO thread_members (thread_id,user_id) VALUES ($1,$2),($1,$3)`, tid, u1, u2); err != nil {
		return "", err
	}
	return tid, nil
}

// listThreadsHandler returns a lightweight summary of direct threads for a user
// GET /v1/threads?user_id=XYZ -> [{thread_id, other_user_id, last_message_at}]
func listThreadsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	userID := r.URL.Query().Get("user_id")
	if userID == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}
	type ThreadSummary struct {
		ThreadID      string    `json:"thread_id"`
		OtherUserID   string    `json:"other_user_id"`
		LastMessageAt time.Time `json:"last_message_at"`
	}
	summaries := []ThreadSummary{}
	if useMemory {
		// Memory mode: derive from memMsgs and memThreads
		memLock.Lock()
		lastByThread := map[string]time.Time{}
		// Build reverse map from thread->pair to identify other user
		pairByThread := map[string][2]string{}
		for pair, tid := range memThreads {
			parts := strings.Split(pair, "|")
			if len(parts) == 2 {
				pairByThread[tid] = [2]string{parts[0], parts[1]}
			}
		}
		for _, m := range memMsgs {
			if m.SentAt.After(lastByThread[m.ThreadID]) {
				lastByThread[m.ThreadID] = m.SentAt
			}
		}
		for tid, pair := range pairByThread {
			if pair[0] == userID || pair[1] == userID {
				other := pair[0]
				if other == userID {
					other = pair[1]
				}
				summaries = append(summaries, ThreadSummary{ThreadID: tid, OtherUserID: other, LastMessageAt: lastByThread[tid]})
			}
		}
		memLock.Unlock()
	} else {
		// SQL mode: find direct threads user participates in, join to other member, get last message time
		q := `SELECT t.id, m2.user_id AS other_user, COALESCE(MAX(msg.sent_at), t.created_at) AS last_message
			  FROM threads t
			  JOIN thread_members m1 ON m1.thread_id=t.id AND m1.user_id=$1
			  JOIN thread_members m2 ON m2.thread_id=t.id AND m2.user_id<>$1
			  LEFT JOIN messages msg ON msg.thread_id=t.id
			  WHERE t.type='direct'
			  GROUP BY t.id, other_user, t.created_at
			  ORDER BY last_message DESC`
		rows, err := db.Query(q, userID)
		if err != nil {
			log.Printf("listThreads query error: %v", err)
			http.Error(w, "internal", http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		for rows.Next() {
			var ts ThreadSummary
			if err := rows.Scan(&ts.ThreadID, &ts.OtherUserID, &ts.LastMessageAt); err != nil {
				log.Printf("scan thread summary: %v", err)
				continue
			}
			summaries = append(summaries, ts)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

// sendDemoHandler and inboxDemoHandler preserve the quick demo endpoints used
// elsewhere in the repository; they operate in-memory even when DB is enabled.
func sendDemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(405)
		return
	}
	var body struct {
		From, To, Content string
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}
	m := Message{ID: uuid.New().String(), SenderID: body.From, ThreadID: uuid.New().String(), Content: body.Content, SentAt: time.Now()}
	memLock.Lock()
	memMsgs = append(memMsgs, m)
	memLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}

func inboxDemoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(405)
		return
	}
	user := r.URL.Query().Get("user_id")
	if user == "" {
		http.Error(w, "missing user_id", http.StatusBadRequest)
		return
	}
	memLock.Lock()
	msgs := []Message{}
	for _, msg := range memMsgs {
		if msg.SenderID == user || strings.Contains(msg.ThreadID, user) {
			msgs = append(msgs, msg)
		}
	}
	memLock.Unlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}
