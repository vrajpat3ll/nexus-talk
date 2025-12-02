package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSendMessageHandler_InvalidMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/v1/messages", nil)
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestSendMessageHandler_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBufferString(`{invalid`))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSendMessageHandler_MissingSenderID(t *testing.T) {
	body := map[string]string{"content": "hello"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestSendMessageHandler_DirectThreadCreation(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "to_id": "u2", "content": "hi"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_ThreadIDProvided(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": "hi"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_SQLInjection(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": "'; DROP TABLE messages;--"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_UnicodeContent(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": "你好世界"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_LongContent(t *testing.T) {
	longContent := ""
	for i := 0; i < 1000; i++ {
		longContent += "a"
	}
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": longContent}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_SpecialCharacters(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": "!@#$%^&*()_+"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_MissingContent(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_EmptyContent(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "t1", "content": ""}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_InvalidThreadID(t *testing.T) {
	body := map[string]string{"sender_id": "u1", "thread_id": "", "content": "hi"}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_MetadataField(t *testing.T) {
	body := map[string]interface{}{"sender_id": "u1", "thread_id": "t1", "content": "hi", "metadata": map[string]interface{}{ "key": "value" }}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}

func TestSendMessageHandler_LargeMetadata(t *testing.T) {
	largeMeta := map[string]interface{}{}
	for i := 0; i < 100; i++ {
		largeMeta[string(rune('a'+i%26))+string(i)] = i
	}
	body := map[string]interface{}{"sender_id": "u1", "thread_id": "t1", "content": "hi", "metadata": largeMeta}
	bodyJSON, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", bytes.NewBuffer(bodyJSON))
	w := httptest.NewRecorder()
	messagesHandler(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Expected %d, got %d", http.StatusOK, w.Code)
	}
}
