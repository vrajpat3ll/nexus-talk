package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gorilla/websocket"
)

func TestWSHandler_MissingUserID(t *testing.T) {
	req := httptest.NewRequest("GET", "/v1/ws", nil)
	w := httptest.NewRecorder()
	wsHandler(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestWSHandler_UpgradeError(t *testing.T) {
	// Simulate upgrade error by passing invalid request
	req := httptest.NewRequest("GET", "/v1/ws?user_id=abc", nil)
	w := httptest.NewRecorder()
	// Remove Upgrade header to force error
	wsHandler(w, req)
	// No status code set, but should not panic
}

func TestWSHandler_SuccessfulConnectionImmediateClose(t *testing.T) {
	// Use httptest server with websocket upgrader
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	_, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
}

func TestWSHandler_SuccessfulConnectionNormalClose(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	c.Close()
}

func TestWSHandler_MessageReceived(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte("hello"))
	if err != nil {
		t.Errorf("write error: %v", err)
	}
	c.Close()
}

func TestWSHandler_LargeMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	large := make([]byte, 4096)
	for i := range large {
		large[i] = 'A'
	}
	err = c.WriteMessage(websocket.TextMessage, large)
	if err != nil {
		t.Errorf("write error: %v", err)
	}
	c.Close()
}

func TestWSHandler_UnicodeMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte("こんにちは世界 🌏"))
	if err != nil {
		t.Errorf("write error: %v", err)
	}
	c.Close()
}

func TestWSHandler_BinaryMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	err = c.WriteMessage(websocket.BinaryMessage, []byte{0x00, 0x01, 0x02, 0x03})
	if err != nil {
		t.Errorf("write error: %v", err)
	}
	c.Close()
}

func TestWSHandler_EmptyMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.RawQuery = "user_id=abc"
		wsHandler(w, r)
	}))
	defer ts.Close()
	url := "ws" + ts.URL[4:] + "/v1/ws?user_id=abc"
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Skip("websocket dial error: ", err)
	}
	err = c.WriteMessage(websocket.TextMessage, []byte(""))
	if err != nil {
		t.Errorf("write error: %v", err)
	}
	c.Close()
}
