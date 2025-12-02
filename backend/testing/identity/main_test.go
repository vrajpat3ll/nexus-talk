package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
)

func setupMockDB(t *testing.T) (sqlmock.Sqlmock, func()) {
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Failed to create mock DB: %v", err)
	}
	db = mockDB
	return mock, func() {
		db.Close()
		db = nil
	}
}

func TestRegisterHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/register", nil)
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestRegisterHandler_InvalidJSON(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	_ = mock
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString(`{invalid`))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterHandler_MissingFields(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	_ = mock
	body := map[string]string{"username": "", "password": "test"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRegisterHandler_DatabaseError(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnError(sql.ErrConnDone)
	body := map[string]string{"username": "testuser", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestRegisterHandler_Success(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "newuser", "password": "securepass", "email": "test@example.com"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
	var response map[string]string
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if _, err := uuid.Parse(response["id"]); err != nil {
		t.Errorf("Invalid UUID: %v", err)
	}
	if response["username"] != "newuser" {
		t.Errorf("Expected username 'newuser', got '%s'", response["username"])
	}
}

func TestRegisterHandler_SQLInjection(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "admin'--", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_UnicodeCharacters(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "用户名", "password": "密码123", "email": "user@test.cn"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_LongInput(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	longUsername := ""
	for i := 0; i < 1000; i++ {
		longUsername += "a"
	}
	body := map[string]string{"username": longUsername, "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_SpecialCharacters(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "user@#$%", "password": "P@ss!w0rd", "email": "test+tag@example.com"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_DuplicateUsername(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnError(sql.ErrConnDone)
	body := map[string]string{"username": "existinguser", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("Expected %d, got %d", http.StatusConflict, w.Code)
	}
}

func TestRegisterHandler_CaseSensitivity(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "TestUser", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_EmptyEmail(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "nomail", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_NullBytes(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	body := map[string]string{"username": "user\x00name", "password": "pass123"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}

func TestRegisterHandler_ExtraFields(t *testing.T) {
	mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectExec("INSERT INTO users").WillReturnResult(sqlmock.NewResult(1, 1))
	payload := map[string]interface{}{
		"username": "user",
		"password": "pass123",
		"email":    "test@test.com",
		"extra":    "ignored",
		"nested":   map[string]string{"key": "value"},
	}
	bodyJSON, _ := json.Marshal(payload)
	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	registerHandler(w, req)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected %d, got %d", http.StatusCreated, w.Code)
	}
}
