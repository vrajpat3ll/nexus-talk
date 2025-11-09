package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Report struct {
	ID            string    `json:"id"`
	ReporterID   string    `json:"reporter_id"`
	TargetID     string    `json:"target_id"`
	ReasonCode   string    `json:"reason_code"`
	Context      bool      `json:"include_context"`
	CreatedAt    time.Time `json:"created_at"`
	Status       string    `json:"status"`
}

type ModerationAction struct {
	ID          string    `json:"id"`
	ReportID    string    `json:"report_id"`
	Action      string    `json:"action"`
	ModeratorID string    `json:"moderator_id"`
	TakenAt     time.Time `json:"taken_at"`
	Comment     string    `json:"comment"`
}

type Block struct {
	BlockerID string `json:"blocker_id"`
	BlockedID string `json:"blocked_id"`
}

type Strike struct {
	ModeratorID string `json:"moderator_id"`
	UserID      string `json:"user_id"`
	Reason      string `json:"reason"`
}

type ModerationService struct {
	db  *sql.DB
	mux *http.ServeMux
}

func NewModerationService() (*ModerationService, error) {
	// TODO: Move connection details to config
	connStr := "postgres://postgres:postgres@localhost:5432/nexustalk?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}
	
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	service := &ModerationService{
		db:  db,
		mux: http.NewServeMux(),
	}

	// Register handlers
	service.mux.HandleFunc("/report", service.CreateReport)
	service.mux.HandleFunc("/block", service.Block)
	service.mux.HandleFunc("/unblock", service.Unblock)
	service.mux.HandleFunc("/strikes", service.IssueStrike)

	return service, nil
}

func (s *ModerationService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *ModerationService) CreateReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var report Report
	if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	report.ID = uuid.New().String()
	report.Status = "new" // Initial state per design doc
	report.CreatedAt = time.Now()

	_, err := s.db.Exec(`
		INSERT INTO reports (id, reporter_id, target_id, reason_code, include_context, status)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, report.ID, report.ReporterID, report.TargetID, report.ReasonCode, report.Context, report.Status)

	if err != nil {
		log.Printf("Failed to create report: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Trigger automated content classification
	go s.classifyReport(report.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(report)
}

func (s *ModerationService) Block(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var block Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
		INSERT INTO block_list (id, blocker_id, blocked_id)
		VALUES ($1, $2, $3)
	`, uuid.New().String(), block.BlockerID, block.BlockedID)

	if err != nil {
		log.Printf("Failed to block user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]bool{"blocked": true})
}

func (s *ModerationService) Unblock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var block Block
	if err := json.NewDecoder(r.Body).Decode(&block); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	_, err := s.db.Exec(`
		DELETE FROM block_list
		WHERE blocker_id = $1 AND blocked_id = $2
	`, block.BlockerID, block.BlockedID)

	if err != nil {
		log.Printf("Failed to unblock user: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *ModerationService) IssueStrike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var strike Strike
	if err := json.NewDecoder(r.Body).Decode(&strike); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Set strike expiry to 30 days from now
	expiry := time.Now().AddDate(0, 0, 30)

	_, err := s.db.Exec(`
		INSERT INTO automod_strikes (id, user_id, reason, expiry)
		VALUES ($1, $2, $3, $4)
	`, uuid.New().String(), strike.UserID, strike.Reason, expiry)

	if err != nil {
		log.Printf("Failed to issue strike: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if user has reached strike limit
	var strikeCount int
	err = s.db.QueryRow(`
		SELECT COUNT(*) FROM automod_strikes
		WHERE user_id = $1 AND active = true
	`, strike.UserID).Scan(&strikeCount)

	if err != nil {
		log.Printf("Failed to count strikes: %v", err)
	} else if strikeCount >= 3 {
		// Implement automated sanctions for users with 3+ strikes
		go s.handleStrikeLimit(strike.UserID)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "strike_issued"})
}

func (s *ModerationService) classifyReport(reportID string) {
	// TODO: Add machine learning classification
	log.Printf("Starting content classification for report: %s", reportID)
	
	// Simulating classification delay
	time.Sleep(2 * time.Second)
	
	// Update report status to "triaged" after classification
	_, err := s.db.Exec(`
		UPDATE reports 
		SET status = 'triaged'
		WHERE id = $1
	`, reportID)
	
	if err != nil {
		log.Printf("Failed to update report status after classification: %v", err)
	}
}

func (s *ModerationService) handleStrikeLimit(userID string) {
	// Implement automated sanctions for repeat offenders
	log.Printf("User %s has reached strike limit, applying sanctions", userID)

	// Record moderation action
	_, err := s.db.Exec(`
		INSERT INTO moderation_actions (id, report_id, action, moderator_id, comment)
		VALUES ($1, NULL, 'auto_suspend', 'system', 'Automated suspension due to strike limit')
	`, uuid.New().String())

	if err != nil {
		log.Printf("Failed to record automated sanction: %v", err)
	}
}

func main() {
	service, err := NewModerationService()
	if err != nil {
		log.Fatalf("Failed to initialize moderation service: %v", err)
	}

	log.Println("Moderation service listening on :8086")
	log.Fatal(http.ListenAndServe(":8086", service))
}