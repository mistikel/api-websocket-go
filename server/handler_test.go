package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func setup(t *testing.T) Handler {
	t.Parallel()
	h := NewHandler()
	return h
}

func TestHealthCheck(t *testing.T) {
	h := setup(t)
	r, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HealthCheck)

	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestHomePage(t *testing.T) {
	h := setup(t)
	r, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.HomePage)

	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestStoreMessage(t *testing.T) {
	h := setup(t)
	body, err := json.Marshal(map[string]string{
		"message": "test message",
	})
	if err != nil {
		t.Fatal(err)
	}
	r, err := http.NewRequest("POST", "/message", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.StoreMessage)

	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}
}

func TestResolveMessage(t *testing.T) {
	h := setup(t)

	r, err := http.NewRequest("GET", "/message", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(h.ResolveMessage)

	handler.ServeHTTP(w, r)

	if status := w.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestSubscribe(t *testing.T) {
	m := setup(t)
	s := httptest.NewServer(http.HandlerFunc(m.SubscribeMessage))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ws.Close()
	if err := ws.WriteMessage(websocket.CloseMessage, nil); err != nil {
		t.Fatal(err)
	}
}
